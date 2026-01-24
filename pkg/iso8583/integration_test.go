package iso8583

import (
	"testing"
)

func TestRoundTrip(t *testing.T) {
	spec := GetDefaultSpec()
	orig := NewMessage()
	orig.MTI = "0200"
	orig.Set(3, "400000")   // Fixed 6
	orig.Set(2, "12345678") // LLVAR

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
	if string(result.Fields[3].Value) != "400000" {
		t.Errorf("Field 3 mismatch: got %s", string(result.Fields[3].Value))
	}
}
