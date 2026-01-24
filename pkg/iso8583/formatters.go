package iso8583

import (
	"fmt"
	"strconv"
)

// Formatter defines how to handle field lengths and data layouts
type Formatter interface {
	Format(val string, length int) ([]byte, error)
	Parse(data []byte, length int) (val []byte, readLen int, err error)
}

// FixedFormatter handles FIXED length fields
type FixedFormatter struct{}

func (f *FixedFormatter) Format(val string, length int) ([]byte, error) {
	// Basic logic: just return bytes (padding logic can be added here)
	return []byte(val), nil
}

func (f *FixedFormatter) Parse(data []byte, length int) ([]byte, int, error) {
	return data[:length], length, nil
}

// LLVarFormatter handles 2-digit length prefixes
type LLVarFormatter struct{}

func (f *LLVarFormatter) Format(val string, length int) ([]byte, error) {
	return []byte(fmt.Sprintf("%02d%s", len(val), val)), nil
}

func (f *LLVarFormatter) Parse(data []byte, length int) ([]byte, int, error) {
	l, _ := strconv.Atoi(string(data[:2]))
	return data[2 : 2+l], 2 + l, nil
}

// LLLVarFormatter handles 3-digit length prefixes
type LLLVarFormatter struct{}

func (f *LLLVarFormatter) Format(val string, length int) ([]byte, error) {
	return []byte(fmt.Sprintf("%03d%s", len(val), val)), nil
}

func (f *LLLVarFormatter) Parse(data []byte, length int) ([]byte, int, error) {
	l, _ := strconv.Atoi(string(data[:3]))
	return data[3 : 3+l], 3 + l, nil
}
