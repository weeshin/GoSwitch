package server

import (
	"GoSwitch/pkg/iso8583"
	"io"
	"log/slog"
	"net"
)

type HandleFunc func(*Context)

type Engine struct {
	Addr           string
	Spec           *iso8583.Spec
	Channel        Channel
	requestHandler HandleFunc
}

func NewEngine(addr string, spec *iso8583.Spec, channel Channel) *Engine {
	return &Engine{
		Addr:    addr,
		Spec:    spec,
		Channel: channel,
	}
}

func (e *Engine) Request(h HandleFunc) {
	e.requestHandler = h
}

func (e *Engine) Start() error {
	ln, err := net.Listen("tcp", e.Addr)
	if err != nil {
		return err
	}
	slog.Info("GoSwitch Framework listening", "addr", e.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Error("Accept error", "error", err)
			continue
		}
		go e.serve(conn)
	}
}

func (e *Engine) serve(conn net.Conn) {
	defer conn.Close()
	slog.Info("New connection", "remote_addr", conn.RemoteAddr())

	for {
		msg, err := e.Channel.Receive(conn)
		if err != nil {
			if err != io.EOF {
				slog.Error("read error", "err", err)
			}
			break
		}
		slog.Info("Received message", "MTI", msg.MTI, "fields", len(msg.Fields))

		// Create Context
		ctx := &Context{
			Request: msg,
			Conn:    conn,
			Channel: e.Channel,
			Spec:    e.Spec,
		}

		// Execute User Logic
		if e.requestHandler != nil {
			go func(c *Context) {
				defer func() {
					if r := recover(); r != nil {
						slog.Error("Panic in request handler", "reason", r)
					}
				}()
				e.requestHandler(c)
			}(ctx)
		}
	}
}
