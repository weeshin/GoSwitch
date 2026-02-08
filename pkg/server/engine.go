package server

import (
	"GoSwitch/pkg/iso8583"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type HandleFunc func(*Context)

type Engine struct {
	Addr             string
	Spec             *iso8583.Spec
	Channel          Channel
	requestHandler   HandleFunc
	slog             *slog.Logger
	Peers            sync.Map
	pendingResponses sync.Map
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
		// Unified loop for incoming connections
		// go e.managePeer(conn, conn.RemoteAddr().String(), true)
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
		ctx := NewContext(msg, sessionChannel, e.Spec, l, e)

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

// Connect adds an Outgoing Peer (Client)
func (e *Engine) Connect(name string, addr string) {
	go func() {
		for {
			conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
			if err != nil {
				e.slog.Error("Failed to connect to peer", "name", name, "addr", addr)
				time.Sleep(5 * time.Second)
				continue
			}

			// Manage the outgoing peer
			e.managePeer(conn, name, false)
			e.slog.Warn("Peer connection lost, retrying...", "name", name)
			time.Sleep(5 * time.Second)
		}
	}()
}

// managePeer is the centralized reader loop for EVERY connection
func (e *Engine) managePeer(conn net.Conn, name string, isIncoming bool) {
	defer conn.Close()

	sessionChannel := e.Channel.Clone(conn)
	e.Peers.Store(name, sessionChannel)
	defer e.Peers.Delete(name)

	e.slog.Info("Peer active", "name", name, "incoming", isIncoming)

	for {
		msg, err := sessionChannel.Receive(conn)
		e.slog.Info("listening for messages...")
		if err != nil {
			e.slog.Error("read error", "err", err, "peer", name)
			break
		}

		// 1. Check if this is a response to something we sent (Correlation)
		e.slog.Info(fmt.Sprintf("Incoming: %s", msg.LogString()))
		stan := msg.Get(11)
		mti := msg.MTI
		ticket := e.createTicket(mti, stan)

		if val, ok := e.pendingResponses.Load(ticket); ok {
			respChan := val.(chan *iso8583.Message)
			respChan <- msg
			continue
		}

	}
}

func (e *Engine) SendAndReceive(peerName string, req *iso8583.Message, timeout time.Duration) (*iso8583.Message, error) {
	// 1. Find the target connection
	val, ok := e.Peers.Load(peerName)
	if !ok {
		return nil, fmt.Errorf("session not found for address: %s", peerName)
	}
	sessionChannel := val.(Channel)

	// 2. Setup correlation (STAN)
	stan := req.Get(11)

	respMTI := predictResponseMTI(req.MTI)
	ticket := e.createTicket(respMTI, stan)

	respChan := make(chan *iso8583.Message, 1)
	e.pendingResponses.Store(ticket, respChan)
	defer e.pendingResponses.Delete(ticket)

	// 3. Send
	if err := sessionChannel.Send(req); err != nil {
		return nil, err
	}

	// 4. Wait
	select {
	case resp := <-respChan:
		return resp, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for STAN %s", ticket)
	}
}

func predictResponseMTI(requestMTI string) string {
	if len(requestMTI) != 4 {
		return ""
	}

	// Convert to byte slice for easy manipulation
	mti := []byte(requestMTI)

	// ISO 8583 Logic:
	// x0xx -> x1xx (Request -> Response)
	// x2xx -> x3xx (Advice -> Advice Response)
	// x4xx -> x5xx (Reversal -> Reversal Response)

	if mti[2] == '0' || mti[2] == '2' || mti[2] == '4' {
		mti[2] = mti[2] + 1
	} else {
		// Fallback: just force it to '1' as you did,
		// which covers the most common 0100/0800 cases.
		mti[2] = '1'
	}

	// Note: Some legacy hosts return 0110 for 0101 (Repeat).
	// We force the 4th digit to '0' to match the response standard.
	mti[3] = '0'

	return string(mti)
}

func (e *Engine) createTicket(mti string, stan string) string {
	// Trim spaces and ensure it's treated consistently
	cleanStan := fmt.Sprintf("%06s", strings.TrimSpace(stan))
	// Take only the last 6 if it's longer
	if len(cleanStan) > 6 {
		cleanStan = cleanStan[len(cleanStan)-6:]
	}
	return fmt.Sprintf("%s_%s", mti, cleanStan)
}
