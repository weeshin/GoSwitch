package iso8583

import (
	"testing"
)

func TestPadding(t *testing.T) {
	// Test Numeric Padding
	gotN := PadLeft("500", 6, "0")
	if gotN != "000500" {
		t.Errorf("Numeric padding failed: got %s", gotN)
	}

	// Test Alpha Padding
	gotA := PadRight("VISA", 10, " ")
	if gotA != "VISA      " {
		t.Errorf("Alpha padding failed: got '%s'", gotA)
	}
}
