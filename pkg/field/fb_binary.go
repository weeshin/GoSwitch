package field

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type FBBinary struct{}

func (f *FBBinary) Pack(val string, length int) ([]byte, error) {
	// The input 'val' is expected to be a Hex string representation of the bytes
	b, err := hex.DecodeString(val)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string for FB_Binary: %v", err)
	}

	// In raw binary, the byte length is usually half the hex string length
	expectedByteLen := length / 2
	if len(b) > expectedByteLen {
		b = b[:expectedByteLen]
	}

	// Pad with null bytes if necessary
	if len(b) < expectedByteLen {
		padded := make([]byte, expectedByteLen)
		copy(padded, b)
		return padded, nil
	}

	return b, nil
}

func (f *FBBinary) Unpack(data []byte, length int) (string, int, error) {
	byteLen := length / 2
	if len(data) < byteLen {
		return "", 0, fmt.Errorf("insufficient data for FB_Binary")
	}

	rawBytes := data[:byteLen]
	return strings.ToUpper(hex.EncodeToString(rawBytes)), byteLen, nil
}
