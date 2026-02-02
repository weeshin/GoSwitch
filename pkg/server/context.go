package server

import (
	"GoSwitch/pkg/iso8583"
	"fmt"
	"log/slog"
)

type Context struct {
	Request *iso8583.Message
	Channel Channel
	Spec    *iso8583.Spec
}

func NewContext(request *iso8583.Message, channel Channel, spec *iso8583.Spec) *Context {
	return &Context{
		Request: request,
		Channel: channel,
		Spec:    spec,
	}
}

// Send packs the message and sends it back using the configured channel
func (c *Context) Send(msg *iso8583.Message) error {
	slog.Info(fmt.Sprintf("Outgoing: %s", msg.LogString()))
	return c.Channel.Send(msg)
}
