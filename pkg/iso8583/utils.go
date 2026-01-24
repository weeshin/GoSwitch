package iso8583

import (
	"strings"
)

// PadLeft adds '0' or spaces to the left (Numeric)
func PadLeft(value string, length int, padChar string) string {
	if len(value) >= length {
		return value[:length]
	}
	return strings.Repeat(padChar, length-len(value)) + value
}

// PadRight adds spaces to the right (Alpha/Alphanumeric)
func PadRight(value string, length int, padChar string) string {
	if len(value) >= length {
		return value[:length]
	}
	return value + strings.Repeat(padChar, length-len(value))
}
