package field

import (
	"fmt"
)

type FBLLChar struct{}

func (f *FBLLChar) Pack(val string, length int) ([]byte, error) {
	dataLen := len(val)
	if dataLen > length {
		return nil, fmt.Errorf("field length %d exceeds max %d", dataLen, length)
	}

	// 1-byte BCD header (e.g., length 12 -> 0x12)
	header := byte(((dataLen / 10) << 4) | (dataLen % 10))

	return append([]byte{header}, []byte(val)...), nil
}

func (f *FBLLChar) Unpack(data []byte, length int) (string, int, error) {
	if len(data) < 1 {
		return "", 0, fmt.Errorf("insufficient data for BCD LL header")
	}

	// Decode 1-byte BCD header
	dataLen := int(data[0]>>4)*10 + int(data[0]&0x0F)

	if len(data) < 1+dataLen {
		return "", 0, fmt.Errorf("insufficient data for FB_LL content")
	}

	return string(data[1 : 1+dataLen]), 1 + dataLen, nil
}
