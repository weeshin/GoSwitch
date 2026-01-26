package field

import (
	"fmt"
)

type FBBitmap struct{}

func (b *FBBitmap) Pack(fields map[int]bool) ([]byte, error) {
	// Determine size (8 or 16 bytes)
	size := 8
	for f := range fields {
		if f > 64 {
			size = 16
			break
		}
	}

	res := make([]byte, size)
	if size == 16 {
		res[0] |= 0x80 // Set Bit 1 for secondary
	}

	for f, present := range fields {
		if !present || f <= 1 || f > 128 {
			continue
		}
		byteIdx := (f - 1) / 8
		bitIdx := uint(7 - ((f - 1) % 8))
		res[byteIdx] |= (1 << bitIdx)
	}
	return res, nil
}

func (b *FBBitmap) Unpack(data []byte) (map[int]bool, int, error) {
	if len(data) < 8 {
		return nil, 0, fmt.Errorf("data too short for primary bitmap")
	}

	fields := make(map[int]bool)
	readLen := 8

	// Logic to check if we need to read 8 or 16 bytes
	hasSecondary := (data[0] & 0x80) != 0
	if hasSecondary {
		readLen = 16
	}

	if len(data) < readLen {
		return nil, 0, fmt.Errorf("data too short for full bitmap")
	}

	for i := 0; i < readLen; i++ {
		for bit := 0; bit < 8; bit++ {
			if (data[i] & (1 << uint(7-bit))) != 0 {
				fieldNum := i*8 + bit + 1
				fields[fieldNum] = true
			}
		}
	}

	return fields, readLen, nil
}
