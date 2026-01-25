package iso8583

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
)

// Field represents a single ISO 8583 data element
type Field struct {
	Value []byte
}

// Message is the main ISO 8583 container
type Message struct {
	MTI    string
	Fields map[int]*Field
}

func NewMessage() *Message {
	return &Message{
		Fields: make(map[int]*Field),
	}
}

// Set adds a field to the message
func (m *Message) Set(fieldNum int, value string) {
	m.Fields[fieldNum] = &Field{Value: []byte(value)}
}

// GenerateBitmap constructs the binary bitmap (8 or 16 bytes)
func (m *Message) GenerateBitmap() ([]byte, error) {
	maxField := 0
	for fieldNum := range m.Fields {
		if fieldNum > maxField {
			maxField = fieldNum
		}
	}

	// Determine if we need a secondary bitmap (fields 65-128)
	size := 8
	if maxField > 64 {
		size = 16
	}

	bitmap := make([]byte, size)

	// If expanding to 16 bytes, the first bit (Field 1) MUST be set
	if size == 16 {
		bitmap[0] |= 0x80
	}

	for fieldNum := range m.Fields {
		if fieldNum < 1 {
			continue
		}
		// Assuming we support up to 128 for now
		if fieldNum > 128 {
			continue
		}

		byteIdx := (fieldNum - 1) / 8
		bitIdx := uint(7 - ((fieldNum - 1) % 8))

		bitmap[byteIdx] |= (1 << bitIdx)
	}

	return bitmap, nil
}

// BitmapHex returns the bitmap as a Hexadecimal string (e.g., "4210001100000000")
func (m *Message) BitmapHex() (string, error) {
	b, err := m.GenerateBitmap()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (m *Message) Pack(spec Spec) ([]byte, error) {
	var buf bytes.Buffer

	// 1. Write MTI (e.g., "0200")
	log.Printf("[DEBUG] Packing MTI: %s", m.MTI)
	buf.WriteString(m.MTI)

	// 2. Generate and Write Bitmap
	bitmap, err := m.GenerateBitmap()
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Generated Bitmap: %x", bitmap)
	buf.Write(bitmap)

	// 3. Loop through fields 2 to 128
	for i := 2; i <= 128; i++ {
		field, exists := m.Fields[i]
		if !exists {
			continue
		}

		fieldSpec, defined := spec[i]
		if !defined {
			return nil, fmt.Errorf("field %d present in message but not in spec", i)
		}

		// Use the Formatter from the Spec to format the field
		formatVal, err := fieldSpec.Formatter.Format(string(field.Value), fieldSpec.Length)
		if err != nil {
			return nil, fmt.Errorf("error formatting field %d: %v", i, err)
		}

		log.Printf("[DEBUG] Packing Field %d: %s", i, formatVal)
		buf.Write(formatVal)
	}

	return buf.Bytes(), nil
}

func (m *Message) Unpack(data []byte, spec Spec) error {
	// 1. Extract MTI (First 4 bytes)
	if len(data) < 2 {
		return fmt.Errorf("data too short for MTI")
	}
	m.MTI = string(data[:2])
	log.Printf("[DEBUG] Unpacking MTI: %s", m.MTI)
	offset := 2

	// 2. Extract Primary Bitmap (Next 8 bytes)
	if len(data) < offset+8 {
		return fmt.Errorf("data too short for Bitmap")
	}
	primaryBitmap := data[offset : offset+8]
	offset += 8

	// Check for Secondary Bitmap (Bit 1 of Primary)
	hasSecondary := (primaryBitmap[0] & 0x80) != 0
	fullBitmap := primaryBitmap
	if hasSecondary {
		fullBitmap = data[offset-8 : offset+8] // Capture 16 bytes
		offset += 8
	}
	log.Printf("[DEBUG] Unpacked Bitmap: %x", fullBitmap)

	// 3. Extract Fields based on Bitmap
	// We start from Field 2 (Bit 1 is the secondary bitmap flag)
	for i := 2; i <= (len(fullBitmap) * 8); i++ {
		byteIdx := (i - 1) / 8
		bitIdx := uint(7 - ((i - 1) % 8))

		// Check if bit is set
		if (fullBitmap[byteIdx] & (1 << bitIdx)) != 0 {
			fSpec, defined := spec[i]
			if !defined {
				return fmt.Errorf("field %d found in bitmap but not in spec", i)
			}

			var fieldVal []byte

			fieldVal, readLen, err := fSpec.Formatter.Parse(data[offset:], fSpec.Length)
			if err != nil {
				return fmt.Errorf("error parsing field %d: %v", i, err)
			}
			log.Printf("[DEBUG] Unpacked Field %d: %s", i, string(fieldVal))
			offset += readLen

			m.Fields[i] = &Field{Value: fieldVal}
		}
	}

	return nil
}
