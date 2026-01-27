package field

import (
	"fmt"
	"strings"
)

type FABinary struct{}

func (f *FABinary) Pack(val string, length int) ([]byte, error) {
	// In ASCII mode, binary data is usually passed as a hex string
	// Ensure it is uppercase and padded/truncated to fixed length
	if len(val) > length {
		val = val[:length]
	}
	padded := val + strings.Repeat("0", length-len(val))
	return []byte(strings.ToUpper(padded)), nil
}

func (f *FABinary) Unpack(data []byte, length int) (string, int, error) {
	if len(data) < length {
		return "", 0, fmt.Errorf("insufficient data for FA_Binary")
	}
	// Return the hex string
	return string(data[:length]), length, nil
}
