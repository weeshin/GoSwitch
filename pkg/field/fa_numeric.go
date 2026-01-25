package field

import (
	"fmt"
	"strings"
)

type FANumeric struct{}

func (f *FANumeric) Pack(val string, length int) ([]byte, error) {
	if len(val) > length {
		val = val[:length] // Or return error
	}
	// Pad left with '0'
	padded := strings.Repeat("0", length-len(val)) + val
	return []byte(padded), nil
}

func (f *FANumeric) Unpack(data []byte, length int) (string, int, error) {
	if len(data) < length {
		return "", 0, fmt.Errorf("insufficient data for FA_Numeric")
	}
	return string(data[:length]), length, nil
}
