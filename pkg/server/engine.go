package server

import (
	"GoSwitch/pkg/iso8583"
	"encoding/hex"
	"io"
	"log/slog"
	"net"
)

type HandleFunc func(*iso8583.Context)

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
		rawData, err := e.Channel.Receive(conn)
		if err != nil {
			if err != io.EOF {
				slog.Error("read error", "err", err)
			}
			break
		}
		slog.Info("Raw", "hex", hex.EncodeToString(rawData))

		msg := iso8583.NewMessage()
		if err := msg.Unpack(rawData, e.Spec); err != nil {
			slog.Error("Unpack Error", "error", err)
			continue
		}

		// Create Context
		ctx := &iso8583.Context{
			Request: msg,
			Conn:    conn,
			Spec:    e.Spec,
		}

		// Execute User Logic
		if e.requestHandler != nil {
			e.requestHandler(ctx)
		}
	}
}
