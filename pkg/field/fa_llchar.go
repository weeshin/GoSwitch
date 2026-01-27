package field

import (
	"fmt"
	"strconv"
)

type FALLChar struct{}

func (f *FALLChar) Pack(val string, length int) ([]byte, error) {
	dataLen := len(val)
	if dataLen > length {
		return nil, fmt.Errorf("field length %d exceeds max %d", dataLen, length)
	}

	// 2-byte ASCII header (e.g., "05")
	header := fmt.Sprintf("%02d", dataLen)
	return append([]byte(header), []byte(val)...), nil
}

func (f *FALLChar) Unpack(data []byte, length int) (string, int, error) {
	if len(data) < 2 {
		return "", 0, fmt.Errorf("insufficient data for LL header")
	}

	// Parse 2-byte ASCII length
	dataLen, err := strconv.Atoi(string(data[:2]))
	if err != nil {
		return "", 0, fmt.Errorf("invalid LL header: %v", err)
	}

	if len(data) < 2+dataLen {
		return "", 0, fmt.Errorf("insufficient data for FA_LL content")
	}

	return string(data[2 : 2+dataLen]), 2 + dataLen, nil
}
