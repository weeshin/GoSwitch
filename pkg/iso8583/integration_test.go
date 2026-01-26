package iso8583

import (
	"GoSwitch/pkg/field"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	spec := &Spec{
		MTIEncoder:    &field.FBNumeric{},
		BitmapEncoder: &field.FBBitmap{},
	}
	spec.Fields = map[int]FieldSpec{
		2: {
			Length:      16,
			Description: "Primary Account Number",
			Encoder:     &field.FBLLNumeric{},
		},
		3: {
			Length:      6,
			Description: "Processing Code",
			Encoder:     &field.FBNumeric{},
		},
	}
	orig := NewMessage()
	orig.MTI = "0200"
	orig.Set(2, "12345678") // LLVAR
	orig.Set(3, "400000")   // Fixed 6

	// Pack it
	packed, err := orig.Pack(spec)
	if err != nil {
		t.Fatalf("Pack failed: %v", err)
	}

	// Unpack into a new object
	result := NewMessage()
	err = result.Unpack(packed, spec)
	if err != nil {
		t.Fatalf("Unpack failed: %v", err)
	}

	// Assertions
	if result.MTI != orig.MTI {
		t.Errorf("MTI mismatch: got %s", result.MTI)
	}
	if string(result.Fields[3].Value) != "300000" {
		t.Errorf("Field 3 mismatch: got %s", string(result.Fields[3].Value))
	}
}
