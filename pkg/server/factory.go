package server

import (
	"GoSwitch/pkg/iso8583"
	"fmt"
	"net"
)

// Constructor defines a function signature that creates a Channel
type Constructor func(conn net.Conn, spec *iso8583.Spec) Channel

var registry = make(map[string]Constructor)

// Register allows you to add new channel types from anywhere in your app
func Register(name string, fn Constructor) {
	registry[name] = fn
}

// NewChannel creates the specific channel based on the config string
func NewChannel(name string, conn net.Conn, spec *iso8583.Spec) (Channel, error) {
	constructor, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown channel type: %s", name)
	}
	return constructor(conn, spec), nil
}
