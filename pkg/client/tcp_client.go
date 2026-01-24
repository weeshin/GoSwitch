package client

import (
	"GoSwitch/pkg/config"
	"GoSwitch/pkg/iso8583"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

type IsoClient struct {
	Config           config.ChannelConfig
	Spec             iso8583.Spec
	conn             net.Conn
	pendingResponses sync.Map // Map[string]chan *iso8583.Message
}

func New(cfg config.ChannelConfig, spec iso8583.Spec) *IsoClient {
	return &IsoClient{
		Config: cfg,
		Spec:   spec,
	}
}

// Connect handles the TCP handshake and keeps retrying if it fails
func (c *IsoClient) Connect() error {
	address := fmt.Sprintf("%s:%d", c.Config.IP, c.Config.Port)

	for {
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)
		if err != nil {
			fmt.Printf("Failed to connect to %s: %v. Retrying in %ds...\n",
				c.Config.Name, err, c.Config.ReconnectInterval)
			time.Sleep(time.Duration(c.Config.ReconnectInterval) * time.Second)
			continue
		}
		c.conn = conn
		fmt.Printf("Connected to channel: %s\n", c.Config.Name)
		return nil
	}
}

func (c *IsoClient) Send(msg *iso8583.Message) error {
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}

	// 1. Pack the ISO message into bytes
	packed, err := msg.Pack(c.Spec)
	if err != nil {
		return err
	}

	// 2. Prepare the 2-byte Length Header
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(len(packed)))

	// 3. Write [Header + Body] to the wire
	_, err = c.conn.Write(append(header, packed...))
	return err
}

// SendAndReceive sends a request and waits for the specific response
func (c *IsoClient) SendAndReceive(req *iso8583.Message, timeout time.Duration) (*iso8583.Message, error) {
	// 1. Extract the Correlation Key (STAN - Field 11)
	stan := string(req.Fields[11].Value)

	// 2. Create a channel to receive the response
	respChan := make(chan *iso8583.Message, 1)
	c.pendingResponses.Store(stan, respChan)

	// Cleanup: Ensure the entry is removed from the map
	defer c.pendingResponses.Delete(stan)

	// 3. Send the message
	if err := c.Send(req); err != nil {
		return nil, err
	}

	// 4. Wait for the response or timeout
	select {
	case resp := <-respChan:
		return resp, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for response for STAN %s", stan)
	}
}
