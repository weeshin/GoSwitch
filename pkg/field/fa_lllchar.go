package field

import (
	"fmt"
	"strconv"
)

type FALLLChar struct{}

func (f *FALLLChar) Pack(val string, length int) ([]byte, error) {
	dataLen := len(val)
	if dataLen > length {
		return nil, fmt.Errorf("field length %d exceeds max %d", dataLen, length)
	}

	// 3-byte ASCII header (e.g., "045" for 45 bytes)
	header := fmt.Sprintf("%03d", dataLen)
	return append([]byte(header), []byte(val)...), nil
}

func (f *FALLLChar) Unpack(data []byte, length int) (string, int, error) {
	if len(data) < 3 {
		return "", 0, fmt.Errorf("insufficient data for LLL header")
	}

	// Parse the 3-byte ASCII header
	dataLen, err := strconv.Atoi(string(data[:3]))
	if err != nil {
		return "", 0, fmt.Errorf("invalid LLL header: %v", err)
	}

	if len(data) < 3+dataLen {
		return "", 0, fmt.Errorf("insufficient data for LLL content: need %d, got %d", 3+dataLen, len(data))
	}

	return string(data[3 : 3+dataLen]), 3 + dataLen, nil
}
