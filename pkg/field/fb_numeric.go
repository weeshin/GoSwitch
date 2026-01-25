package field

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type FBNumeric struct{}

func (f *FBNumeric) Pack(val string, length int) ([]byte, error) {
	// 1. Pad to the required length
	padded := strings.Repeat("0", length-len(val)) + val

	// 2. BCD requires an even number of digits to form full bytes
	if len(padded)%2 != 0 {
		padded = "0" + padded
	}

	return hex.DecodeString(padded)
}

func (f *FBNumeric) Unpack(data []byte, length int) (string, int, error) {
	// Binary length is half the field length (rounded up)
	byteLen := (length + 1) / 2
	if len(data) < byteLen {
		return "", 0, fmt.Errorf("insufficient data for FB_Numeric")
	}

	raw := data[:byteLen]
	res := hex.EncodeToString(raw)

	// If length was odd, we might have a leading padding zero from packing
	if len(res) > length {
		res = res[len(res)-length:]
	}

	return res, byteLen, nil
}
