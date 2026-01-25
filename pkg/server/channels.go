package server

import (
	"encoding/binary"
	"io"
)

// Channel defines how to read/write the message length header
type Channel interface {
	ReadLength(r io.Reader) (int, error)
	WriteLength(w io.Writer, length int) error
}

// NACChannel: 2-byte binary length (Big Endian)
type NACChannel struct{}

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

// BCDChannel: 2-byte BCD length (e.g., 123 bytes = 0x01 0x23)
type BCDChannel struct{}

func (b *BCDChannel) ReadLength(r io.Reader) (int, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, err
	}
	// Convert BCD bytes to integer
	return int(header[0]>>4)*1000 + int(header[0]&0x0F)*100 + int(header[1]>>4)*10 + int(header[1]&0x0F), nil
}

func (b *BCDChannel) WriteLength(w io.Writer, length int) error {
	// Simple BCD encoding for 4 digits
	b1 := byte(((length / 1000) << 4) | ((length / 100) % 10))
	b2 := byte((((length / 10) % 10) << 4) | (length % 10))
	_, err := w.Write([]byte{b1, b2})
	return err
}
