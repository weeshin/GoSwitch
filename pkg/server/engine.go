package server

import (
	"GoSwitch/pkg/iso8583"
	"encoding/hex"
	"io"
	"log"
	"net"
)

type HandleFunc func(*iso8583.Context)

type Engine struct {
	Addr           string
	Spec           iso8583.Spec
	Channel        Channel
	requestHandler HandleFunc
}

func NewEngine(addr string, spec iso8583.Spec, channel Channel) *Engine {
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
	log.Printf("GoSwitch Framework listening on %s\n", e.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go e.serve(conn)
	}
}

func (e *Engine) serve(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection from %s", conn.RemoteAddr())

	for {
		// 1. Read Length using the selected Channel strategy
		msgLen, err := e.Channel.ReadLength(conn)
		log.Printf("Message Length: %d", msgLen)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error (length): %v", err)
			}
			return
		}

		// header := make([]byte, 2)
		// if _, err := io.ReadFull(conn, header); err != nil {
		// 	return
		// }

		body := make([]byte, msgLen)
		if _, err := io.ReadFull(conn, body); err != nil {
			log.Printf("Read error (body): %v", err)
			return
		}
		log.Printf("Body : %v", hex.EncodeToString(body))

		msg := iso8583.NewMessage()
		if err := msg.Unpack(body, e.Spec); err != nil {
			log.Printf("Unpack Error: %v\n", err)
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
