package field

import (
	"fmt"
	"strings"
)

type FABitmap struct{}

func (f *FABitmap) Pack(val string, length int) ([]byte, error) {
	// val is expected to be the hex string representation
	// length here is usually 16 (primary) or 32 (secondary)
	if len(val) < length {
		val = val + strings.Repeat("0", length-len(val))
	}
	return []byte(strings.ToUpper(val)), nil
}

func (f *FABitmap) Unpack(data []byte, length int) (string, int, error) {
	// For ASCII, we read 'length' characters (e.g., 16 or 32)
	if len(data) < length {
		return "", 0, fmt.Errorf("insufficient data for FA_Bitmap")
	}
	return string(data[:length]), length, nil
}
