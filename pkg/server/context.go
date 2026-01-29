package server

import (
	"GoSwitch/pkg/iso8583"
	"net"
)

type Context struct {
	Request *iso8583.Message
	Conn    net.Conn
	Channel Channel
	Spec    *iso8583.Spec
}

// Respond packs the message and sends it back with the 2-byte length header
// func (c *Context) Respond(m *iso8583.Message) error {
// 	// Implementation relying on manual packing.
// 	// Ideally usage of Channel.Send() is preferred if available.
// 	return nil
// }
