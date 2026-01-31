package server

import (
	"GoSwitch/pkg/iso8583"
	"encoding/binary"
	"io"
	"net"
)

type NACHeader struct{}

func (h *NACHeader) ReadLength(r io.Reader) (int, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint16(header)), nil
}

func (h *NACHeader) WriteLength(w io.Writer, length int) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(length))
	_, err := w.Write(header)
	return err
}

// NACChannel: 2-byte binary length (Big Endian)
type NACChannel struct {
	// Some hosts expect a specific TPDU (e.g., 6000000000)
	*BaseChannel
}

func NewNACChannel(conn net.Conn, spec *iso8583.Spec) Channel {
	return &NACChannel{
		BaseChannel: &BaseChannel{
			Conn:    conn,
			Spec:    spec,
			Handler: &NACHeader{},
		},
	}
}

func init() {
	Register("NAC", func(conn net.Conn, spec *iso8583.Spec) Channel {
		return NewNACChannel(conn, spec)
	})
}

func (n *NACChannel) ReadLength(r io.Reader) (int, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint16(header)), nil
}

func (n *NACChannel) WriteLength(w io.Writer, length int) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(length))
	_, err := w.Write(header)
	return err
}

func (n *NACChannel) Receive(r io.Reader) (*iso8583.Message, error) {
	// 1. Read TCP Length
	header := make([]byte, 2)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}
	totalLen := int(binary.BigEndian.Uint16(header))

	// 2. Read full packet (TPDU + Data)
	payload := make([]byte, totalLen)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	isoData := payload
	var msgTPDU []byte

	// 3. Handle TPDU if present
	if len(n.Header) > 0 && totalLen >= 5 {
		msgTPDU = payload[:5]
		isoData = payload[5:]
	}

	// 4. Unpack using the injected Spec
	msg := iso8583.NewMessage()
	if err := msg.Unpack(isoData, n.Spec); err != nil {
		return nil, err
	}

	// Optional: Store the received TPDU in the message
	// so we can swap it during Send
	msg.SetHeader(msgTPDU)

	return msg, nil
}

func (n *NACChannel) Send(w io.Writer, msg *iso8583.Message) error {
	// 1. Pack the ISO message to bytes
	isoBytes, err := msg.Pack(n.Spec)
	if err != nil {
		return err
	}

	// 2. Handle TPDU (Swap Source/Dest if necessary)
	finalPayload := isoBytes
	if len(n.Header) > 0 {
		tpdu := msg.GetHeader()
		if len(tpdu) == 5 {
			// Simple Swap: Swap bytes 1-2 with 3-4
			swapped := []byte{tpdu[0], tpdu[3], tpdu[4], tpdu[1], tpdu[2]}
			finalPayload = append(swapped, isoBytes...)
		} else {
			finalPayload = append(n.Header, isoBytes...)
		}
	}

	// 3. Write TCP Length + Payload
	lenHeader := make([]byte, 2)
	binary.BigEndian.PutUint16(lenHeader, uint16(len(finalPayload)))

	// if _, err := w.Write(lenHeader); err != nil {
	// 	return err
	// }
	// _, err = w.Write(finalPayload)
	if _, err := w.Write(lenHeader); err != nil {
		return err
	}
	_, err = w.Write(finalPayload)
	return err
}
