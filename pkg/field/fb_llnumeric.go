package field

import (
	"encoding/hex"
	"fmt"
)

type FBLLNumeric struct{}

func (f *FBLLNumeric) Pack(val string, length int) ([]byte, error) {
	dataLen := len(val)
	if dataLen > length {
		return nil, fmt.Errorf("field length %d exceeds max %d", dataLen, length)
	}

	// 1. Pack Length into 1 BCD byte (e.g., len 12 -> 0x12)
	header := byte(((dataLen / 10) << 4) | (dataLen % 10))

	// 2. Pack Data into BCD
	padded := val
	if len(padded)%2 != 0 {
		padded = "0" + padded // Ensure even for BCD packing
	}
	dataBytes, _ := hex.DecodeString(padded)

	return append([]byte{header}, dataBytes...), nil
}

func (f *FBLLNumeric) Unpack(data []byte, length int) (string, int, error) {
	if len(data) < 1 {
		return "", 0, fmt.Errorf("insufficient data for BCD LL header")
	}

	// Read 1-byte BCD header (e.g., 0x12 -> int 12)
	dataLen := int(data[0]>>4)*10 + int(data[0]&0x0F)

	// Calculate expected BCD data length
	byteLen := (dataLen + 1) / 2
	if len(data) < 1+byteLen {
		return "", 0, fmt.Errorf("insufficient data for FB_LL content")
	}

	raw := data[1 : 1+byteLen]
	res := hex.EncodeToString(raw)

	// Trim leading zero if dataLen was odd
	if len(res) > dataLen {
		res = res[len(res)-dataLen:]
	}

	return res, 1 + byteLen, nil
}
