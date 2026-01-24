package server

import (
	"GoSwitch/pkg/iso8583"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type IsoServer struct {
	Addr string
	Spec iso8583.Spec
}

func New(addr string, spec iso8583.Spec) *IsoServer {
	return &IsoServer{
		Addr: addr,
		Spec: spec,
	}
}

func (s *IsoServer) Start() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	fmt.Printf("GoSwitch Server started on %s\n", s.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.handle(conn)
	}
}

func (s *IsoServer) handle(conn net.Conn) {
	defer conn.Close()

	for {
		// 1. Read 2-byte length header
		header := make([]byte, 2)
		if _, err := io.ReadFull(conn, header); err != nil {
			return
		}
		msgLen := binary.BigEndian.Uint16(header)

		// 2. Read body
		body := make([]byte, msgLen)
		if _, err := io.ReadFull(conn, body); err != nil {
			return
		}

		// 3. Unpack using the core library
		msg := iso8583.NewMessage()
		if err := msg.Unpack(body, s.Spec); err != nil {
			fmt.Printf("Unpack Error: %v\n", err)
			continue
		}

		fmt.Printf("Incoming: MTI %s, STAN %s\n", msg.MTI, string(msg.Fields[11].Value))
	}
}
