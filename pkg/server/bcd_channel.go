package server

import (
	"GoSwitch/pkg/iso8583"
	"io"
	"net"
)

// BCDChannel: 2-byte BCD length (e.g., 123 bytes = 0x01 0x23)
type BCDChannel struct {
	Spec *iso8583.Spec
}

func NewBCDChannel(conn net.Conn, spec *iso8583.Spec) Channel {
	return &BCDChannel{Spec: spec}
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
	length, err := b.ReadLength(r)
	if err != nil {
		return nil, err
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	msg := iso8583.NewMessage()
	if err := msg.Unpack(payload, b.Spec); err != nil {
		return nil, err
	}
	return msg, nil
}

func (b *BCDChannel) Send(w io.Writer, msg *iso8583.Message) error {
	data, err := msg.Pack(b.Spec)
	if err != nil {
		return err
	}

	if err := b.WriteLength(w, len(data)); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
