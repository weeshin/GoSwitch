package server

import (
	"GoSwitch/pkg/iso8583"
	"io"
	"net"
)

type Channel interface {
	ReadLength(r io.Reader) (int, error)
	WriteLength(w io.Writer, length int) error
	// Receive reads the length and then the full message body
	Receive(r io.Reader) (*iso8583.Message, error)
	// Send writes the length header and then the message body
	Send(msg *iso8583.Message) error
	// Clone creates a new Channel instance with the given connection
	Clone(conn net.Conn) Channel
}

type LengthHandler interface {
	ReadLength(r io.Reader) (int, error)
	WriteLength(w io.Writer, length int) error
}

type BaseChannel struct {
	Conn    net.Conn
	Spec    *iso8583.Spec
	Header  []byte
	Handler LengthHandler
}

func NewBaseChannel(conn net.Conn, spec *iso8583.Spec) *BaseChannel {
	return &BaseChannel{
		Conn: conn,
		Spec: spec,
	}
}

func (b *BaseChannel) ReadLength(r io.Reader) (int, error) {
	return b.Handler.ReadLength(r)
}

func (b *BaseChannel) WriteLength(w io.Writer, length int) error {
	return b.Handler.WriteLength(w, length)
}

func (b *BaseChannel) Receive(r io.Reader) (*iso8583.Message, error) {
	// 1. Read Length via the specialized handler
	length, err := b.Handler.ReadLength(r)
	if err != nil {
		return nil, err
	}

	// 2. Read Body
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	isoData := payload
	var msgHeader []byte
	if len(b.Header) > 0 && length >= 5 {
		msgHeader = payload[:5]
		isoData = payload[5:]
	}

	msg := iso8583.NewMessage()
	if err := msg.Unpack(isoData, b.Spec); err != nil {
		return nil, err
	}
	msg.SetHeader(msgHeader)
	return msg, nil
}

func (b *BaseChannel) Send(msg *iso8583.Message) error {
	isoBytes, err := msg.Pack(b.Spec)
	if err != nil {
		return err
	}
	// ... insert your TPDU swapping logic here ...

	finalPayload := isoBytes

	// Write using the specialized length handler
	if err := b.Handler.WriteLength(b.Conn, len(finalPayload)); err != nil {
		return err
	}
	_, err = b.Conn.Write(finalPayload)
	return err
}

func (b *BaseChannel) Clone(conn net.Conn) Channel {
	if conn == nil {
		// Log error?
		panic("BaseChannel.Clone received nil conn")
	}
	return &BaseChannel{
		Conn:    conn,
		Spec:    b.Spec,
		Header:  b.Header,
		Handler: b.Handler,
	}
}
