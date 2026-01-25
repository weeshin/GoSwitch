package iso8583

import (
	"encoding/binary"
	"net"
)

type Context struct {
	Request *Message
	Conn    net.Conn
	Spec    Spec
}

// Respond packs the message and sends it back with the 2-byte length header
func (c *Context) Respond(m *Message) error {
	packed, err := m.Pack(c.Spec)
	if err != nil {
		return err
	}

	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(len(packed)))

	_, err = c.Conn.Write(append(header, packed...))
	return err
}
