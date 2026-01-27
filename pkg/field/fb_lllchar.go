package field

import (
	"fmt"
)

type FBLLLChar struct{}

func (f *FBLLLChar) Pack(val string, length int) ([]byte, error) {
	dataLen := len(val)
	if dataLen > length {
		return nil, fmt.Errorf("field length %d exceeds max %d", dataLen, length)
	}

	// 2-byte BCD header for 3-4 digits (e.g., 125 -> 0x01, 0x25)
	header := make([]byte, 2)
	header[0] = byte(dataLen / 100)                         // Hundreds
	header[1] = byte(((dataLen/10)%10)<<4 | (dataLen % 10)) // Tens and Units

	return append(header, []byte(val)...), nil
}

func (f *FBLLLChar) Unpack(data []byte, length int) (string, int, error) {
	if len(data) < 2 {
		return "", 0, fmt.Errorf("insufficient data for BCD LLL header")
	}

	// Decode 2-byte BCD header
	hundreds := int(data[0])
	tensUnits := int(data[1])
	dataLen := (hundreds * 100) + (int(tensUnits>>4) * 10) + (tensUnits & 0x0F)

	if len(data) < 2+dataLen {
		return "", 0, fmt.Errorf("insufficient data for FB_LLL content")
	}

	return string(data[2 : 2+dataLen]), 2 + dataLen, nil
}
