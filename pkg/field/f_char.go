package field

import (
	"fmt"
	"strings"
)

type FChar struct{}

// Pack pads the string with spaces on the right to reach the fixed length
func (f *FChar) Pack(val string, length int) ([]byte, error) {
	if len(val) > length {
		// Most hosts truncate, but returning an error is safer for a framework
		val = val[:length]
	}

	// Right pad with spaces (e.g., "ABC" length 5 -> "ABC  ")
	padded := val + strings.Repeat(" ", length-len(val))
	return []byte(padded), nil
}

// Unpack reads exactly 'length' bytes from the data
func (f *FChar) Unpack(data []byte, length int) (string, int, error) {
	if len(data) < length {
		return "", 0, fmt.Errorf("insufficient data for F_Char: need %d, got %d", length, len(data))
	}

	// We return the string including spaces.
	// The user can .TrimSpace() later if they wish.
	return string(data[:length]), length, nil
}
