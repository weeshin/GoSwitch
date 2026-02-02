package server

import (
	"GoSwitch/pkg/iso8583"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

type HandleFunc func(*Context)

type Engine struct {
	Addr           string
	Spec           *iso8583.Spec
	Channel        Channel
	requestHandler HandleFunc
	slog           *slog.Logger
}

func NewEngine(addr string, spec *iso8583.Spec, channel Channel) *Engine {
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				// Custom time format with milliseconds
				return slog.String(slog.TimeKey, a.Value.Time().Format("2006-01-02 15:04:05.000"))
			}
			return a
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	return &Engine{
		Addr:    addr,
		Spec:    spec,
		Channel: channel,
		slog:    logger,
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
	e.slog.Info("GoSwitch Framework listening", "addr", e.Addr)

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

	sessionChannel := e.Channel.Clone(conn)
	l := e.slog.With("remote_addr", conn.RemoteAddr())
	l.Info("New connection")

	for {
		msg, err := sessionChannel.Receive(conn)
		if err != nil {
			if err != io.EOF {
				l.Error("read error", "err", err)
			}
			break
		}
		l.Info(fmt.Sprintf("Incoming: %s", msg.LogString()))
		// Create Context
		ctx := NewContext(msg, sessionChannel, e.Spec, l)

		// Execute User Logic
		if e.requestHandler != nil {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						l.Error("Panic in request handler", "reason", r)
					}
				}()
				e.requestHandler(ctx)
			}()
		}
	}
}
