package server

import (
	"GoSwitch/pkg/iso8583"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// BASE24TCPChannel implements ACI's BASE24 TCP variant
type BASE24TCPChannel struct {
	Conn   net.Conn
	Spec   *iso8583.Spec
	Header []byte // Generic header (e.g., 10-byte ASCII)
}

func NewBASE24TCPChannel(conn net.Conn, spec *iso8583.Spec) Channel {
	return &BASE24TCPChannel{
		Conn: conn,
		Spec: spec,
	}
}

func init() {
	Register("BASE24", func(conn net.Conn, spec *iso8583.Spec) Channel {
		return NewBASE24TCPChannel(conn, spec)
	})
}

func (b *BASE24TCPChannel) ReadLength(r io.Reader) (int, error) {
	header := make([]byte, 2)
	for {
		if _, err := io.ReadFull(r, header); err != nil {
			return 0, err
		}
		length := int(binary.BigEndian.Uint16(header))

		if length == 0 {
			// BASE24 Keep-Alive logic: echo 0 back
			if b.Conn != nil {
				b.Conn.Write(header)
			}
			continue
		}

		// Total length includes the 1-byte trailer, so message length is l-1
		return length - 1, nil
	}
}

func (b *BASE24TCPChannel) WriteLength(w io.Writer, length int) error {
	header := make([]byte, 2)
	// BASE24 length = payload length + 1 (for trailer)
	binary.BigEndian.PutUint16(header, uint16(length+1))
	_, err := w.Write(header)
	return err
}

func (b *BASE24TCPChannel) Receive(r io.Reader) (*iso8583.Message, error) {
	// 1. Read Length (handles the 0-length keep-alive internally)
	msgLen, err := b.ReadLength(r)
	if err != nil {
		return nil, err
	}

	// 2. Read Body
	payload := make([]byte, msgLen)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	// 3. Read Trailer (1 byte, usually 0x03)
	trailer := make([]byte, 1)
	if _, err := io.ReadFull(r, trailer); err != nil {
		return nil, err
	}

	isoData := payload
	var msgHeader []byte
	// 4. Handle Header (General Header, e.g., 10 bytes)
	// If we expect a header, we extract it from the front
	if len(b.Header) > 0 && msgLen >= len(b.Header) {
		hLen := len(b.Header)
		msgHeader = payload[:hLen]
		isoData = payload[hLen:]
	}

	// 4. Unpack
	msg := iso8583.NewMessage()
	// Note: BASE24 often doesn't use a TPDU in this specific TCP variant,
	// but we'll unpack the whole payload.
	if err := msg.Unpack(isoData, b.Spec); err != nil {
		return nil, err
	}
	msg.SetHeader(msgHeader)

	return msg, nil
}

func (b *BASE24TCPChannel) Send(msg *iso8583.Message) error {
	if b.Conn == nil {
		return fmt.Errorf("BASE24TCPChannel.Conn is nil")
	}

	isoBytes, err := msg.Pack(b.Spec)
	if err != nil {
		return err
	}

	// 1. Prepare Header
	h := msg.GetHeader()
	if h == nil {
		h = b.Header
	}

	// 2. Combine Header + ISO Body
	finalPayload := isoBytes
	if h != nil {
		finalPayload = append(h, isoBytes...)
	}

	// 3. Write Length (Total + 1)
	if err := b.WriteLength(b.Conn, len(finalPayload)); err != nil {
		return err
	}

	// 4. Write Data
	if _, err := b.Conn.Write(finalPayload); err != nil {
		return err
	}

	// 5. Write Trailer (ETX 0x03)
	_, err = b.Conn.Write([]byte{0x03})
	return err
}

func (b *BASE24TCPChannel) Clone(conn net.Conn) Channel {
	return &BASE24TCPChannel{
		Conn:   conn,
		Spec:   b.Spec,
		Header: b.Header,
	}
}
