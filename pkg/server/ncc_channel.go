package server

import (
	"GoSwitch/pkg/iso8583"
	"fmt"
	"io"
	"net"
)

// NCCChannel: 4-byte ASCII length (e.g., "0123")
type NCCChannel struct {
	Conn   net.Conn
	Spec   *iso8583.Spec
	Header []byte // Used for TPDU
}

func NewNCCChannel(conn net.Conn, spec *iso8583.Spec) Channel {
	return &NCCChannel{
		Conn: conn,
		Spec: spec,
	}
}

func init() {
	Register("NCC", func(conn net.Conn, spec *iso8583.Spec) Channel {
		return NewNCCChannel(conn, spec)
	})
}

func (n *NCCChannel) ReadLength(r io.Reader) (int, error) {
	b := make([]byte, 2)
	if _, err := io.ReadFull(r, b); err != nil {
		return 0, err
	}
	// BCD to Int: 0x12 0x34 -> 1234
	// (b[0]>>4)*1000 + (b[0]&0x0F)*100 + (b[1]>>4)*10 + (b[1]&0x0F)
	high := int(b[0]>>4)*10 + int(b[0]&0x0F)
	low := int(b[1]>>4)*10 + int(b[1]&0x0F)
	return high*100 + low, nil
}

func (n *NCCChannel) WriteLength(w io.Writer, length int) error {
	// Ensure we only handle up to 9999
	l := length % 10000
	b := make([]byte, 2)
	// Int to BCD: 1234 -> 0x12 0x34
	b[0] = byte(((l / 1000) << 4) | ((l / 100) % 10))
	b[1] = byte((((l / 10) % 10) << 4) | (l % 10))
	_, err := w.Write(b)
	return err
}

func (n *NCCChannel) Receive(r io.Reader) (*iso8583.Message, error) {
	// 1. Read 2-byte ASCII length
	totalLen, err := n.ReadLength(r)
	if err != nil {
		return nil, err
	}

	// 2. Read full packet
	payload := make([]byte, totalLen)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	isoData := payload
	var msgTPDU []byte

	// 3. Handle TPDU (if 5-byte TPDU exists)
	if len(n.Header) > 0 && totalLen >= 5 {
		msgTPDU = payload[:5]
		isoData = payload[5:]
	}

	msg := iso8583.NewMessage()
	if err := msg.Unpack(isoData, n.Spec); err != nil {
		return nil, err
	}

	msg.SetHeader(msgTPDU)
	return msg, nil
}

func (n *NCCChannel) Send(msg *iso8583.Message) error {
	if n.Conn == nil {
		return fmt.Errorf("NCCChannel.Conn is nil")
	}

	isoBytes, err := msg.Pack(n.Spec)
	if err != nil {
		return err
	}

	finalPayload := isoBytes
	if len(n.Header) > 0 {
		tpdu := msg.GetHeader()
		if len(tpdu) == 5 {
			// TPDU Swap Logic
			swapped := []byte{tpdu[0], tpdu[3], tpdu[4], tpdu[1], tpdu[2]}
			finalPayload = append(swapped, isoBytes...)
		} else {
			finalPayload = append(n.Header, isoBytes...)
		}
	}

	// 3. Write 4-byte ASCII Length + Payload
	if err := n.WriteLength(n.Conn, len(finalPayload)); err != nil {
		return err
	}
	_, err = n.Conn.Write(finalPayload)
	return err
}

func (n *NCCChannel) Clone(conn net.Conn) Channel {
	return &NCCChannel{
		Conn:   conn,
		Spec:   n.Spec,
		Header: n.Header,
	}
}
