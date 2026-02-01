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
