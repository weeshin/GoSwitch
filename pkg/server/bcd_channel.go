package server

import (
	"GoSwitch/pkg/iso8583"
	"fmt"
	"io"
	"net"
)

// BCDChannel: 2-byte BCD length (e.g., 123 bytes = 0x01 0x23)
type BCDChannel struct {
	Spec   *iso8583.Spec
	Conn   net.Conn
	Header []byte
}

func NewBCDChannel(conn net.Conn, spec *iso8583.Spec) Channel {
	return &BCDChannel{Spec: spec, Conn: conn}
}

func init() {
	Register("BCD", func(conn net.Conn, spec *iso8583.Spec) Channel {
		return NewBCDChannel(conn, spec)
	})
}

func (b *BCDChannel) ReadLength(r io.Reader) (int, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, err
	}
	// Convert BCD bytes to integer
	return int(header[0]>>4)*1000 + int(header[0]&0x0F)*100 + int(header[1]>>4)*10 + int(header[1]&0x0F), nil
}

func (b *BCDChannel) WriteLength(w io.Writer, length int) error {
	// Simple BCD encoding for 4 digits
	b1 := byte(((length / 1000) << 4) | ((length / 100) % 10))
	b2 := byte((((length / 10) % 10) << 4) | (length % 10))
	_, err := w.Write([]byte{b1, b2})
	return err
}

func (b *BCDChannel) Receive(r io.Reader) (*iso8583.Message, error) {
	// 1. Read TCP Length
	totalLen, err := b.ReadLength(r)
	if err != nil {
		return nil, err
	}

	// 2. Read Message Payload
	payload := make([]byte, totalLen)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	isoData := payload
	var msgTPDU []byte

	// 3. Handle TPDU if present
	if len(b.Header) > 0 && totalLen >= 5 {
		msgTPDU = payload[:5]
		isoData = payload[5:]
	}

	msg := iso8583.NewMessage()
	if err := msg.Unpack(isoData, b.Spec); err != nil {
		return nil, err
	}
	msg.SetHeader(msgTPDU)

	return msg, nil
}

func (b *BCDChannel) Send(msg *iso8583.Message) error {
	if b.Conn == nil {
		return fmt.Errorf("BCDChannel.Conn is nil")
	}

	// 1. Pack the ISO message to bytes
	isoBytes, err := msg.Pack(b.Spec)
	if err != nil {
		return err
	}

	// 2. Handle TPDU (Swap Source/Dest if necessary)
	finalPayload := isoBytes
	if len(b.Header) > 0 {
		tpdu := msg.GetHeader()
		if len(tpdu) == 5 {
			// Simple Swap: Swap bytes 1-2 with 3-4
			swapped := []byte{tpdu[0], tpdu[3], tpdu[4], tpdu[1], tpdu[2]}
			finalPayload = append(swapped, isoBytes...)
		} else {
			finalPayload = append(b.Header, isoBytes...)
		}
	}

	// 3. Write TCP Length + Payload
	if err := b.WriteLength(b.Conn, len(finalPayload)); err != nil {
		return err
	}
	_, err = b.Conn.Write(finalPayload)
	return err
}

func (b *BCDChannel) Clone(conn net.Conn) Channel {
	return &BCDChannel{
		Spec:   b.Spec,
		Conn:   conn,
		Header: b.Header,
	}
}
