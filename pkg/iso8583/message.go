package iso8583

import (
	"bytes"
	"fmt"
	"log/slog"
)

// Field represents a single ISO 8583 data element
type Field struct {
	Value []byte
}

// Message is the main ISO 8583 container
type Message struct {
	MTI    string
	Bitmap []byte
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

// GenerateBitmapHex constructs the binary bitmap (8 or 16 bytes)
func (m *Message) GenerateBitmapHex() ([]byte, error) {
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

func (m *Message) Pack(spec *Spec) ([]byte, error) {
	var buf bytes.Buffer

	// 1. Pack MTI
	// We use the encoder defined in the spec to handle ASCII/Binary MTI
	mtiBytes, err := spec.MTIEncoder.Pack(m.MTI, 4)
	if err != nil {
		return nil, err
	}
	buf.Write(mtiBytes)

	// 2. Pack Bitmap - No more manual bit manipulation here!
	presentFields := make(map[int]bool)
	for k := range m.Fields {
		presentFields[k] = true
	}

	bitmapB, err := spec.BitmapEncoder.Pack(presentFields)
	if err != nil {
		return nil, err
	}
	buf.Write(bitmapB)

	// 3. Pack Fields
	// We loop based on the Spec to ensure we only pack what's allowed
	for i := 2; i <= 128; i++ {
		if fData, ok := m.Fields[i]; ok {
			fSpec, ok := spec.Fields[i]
			if !ok {
				return nil, fmt.Errorf("field %d present in message but not in spec", i)
			}
			packed, err := fSpec.Encoder.Pack(string(fData.Value), fSpec.Length)
			if err != nil {
				return nil, err
			}
			buf.Write(packed)
		}
	}
	return buf.Bytes(), nil
}

func (m *Message) Unpack(data []byte, spec *Spec) error {
	offset := 0

	// 1. Unpack MTI
	mti, readLen, err := spec.MTIEncoder.Unpack(data[offset:], 4)
	if err != nil {
		return err
	}
	m.MTI = mti
	offset += readLen
	slog.Debug("Unpacked MTI", "mti", m.MTI)

	// 2. Unpack Bitmap using the specialized BitMap interface
	// This now returns a map[int]bool directly!
	presentFields, readLen, err := spec.BitmapEncoder.Unpack(data[offset:])
	if err != nil {
		return err
	}
	// Store raw bytes for historical/debug purposes
	m.Bitmap = data[offset : offset+readLen]
	offset += readLen

	slog.Debug("Unpacked Bitmap fields", "fields", presentFields)

	// 3. Extract Fields based on the map returned by the encoder
	// We loop from field 2 up to 128
	for i := 2; i <= 128; i++ {
		// Only unpack if the bitmap says the field is present
		if presentFields[i] {
			fSpec, defined := spec.Fields[i]
			if !defined {
				return fmt.Errorf("field %d found in bitmap but not in spec", i)
			}

			// Unpack the data field
			val, readLen, err := fSpec.Encoder.Unpack(data[offset:], fSpec.Length)
			if err != nil {
				return fmt.Errorf("error unpacking field %d: %v", i, err)
			}

			slog.Debug("Unpacked Field", "field", i, "value", val)

			// Store the value in the message
			m.Fields[i] = &Field{Value: []byte(val)}
			offset += readLen
		}
	}

	return nil
}
