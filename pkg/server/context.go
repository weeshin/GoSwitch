package server

import (
	"GoSwitch/pkg/iso8583"
	"net"
)

type Context struct {
	Request *iso8583.Message
	conn    net.Conn
	Channel Channel
	Spec    *iso8583.Spec
}

func NewContext(conn net.Conn, request *iso8583.Message, channel Channel, spec *iso8583.Spec) *Context {
	return &Context{
		Request: request,
		conn:    conn,
		Channel: channel,
		Spec:    spec,
	}
}

// Send packs the message and sends it back using the configured channel
func (c *Context) Send(msg *iso8583.Message) error {
	return c.Channel.Send(c.conn, msg)
}
