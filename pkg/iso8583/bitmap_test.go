package iso8583

import (
	"testing"
)

func TestBitmapHex(t *testing.T) {
	tests := []struct {
		name     string
		fields   []int
		expected string
	}{
		{
			name:     "Empty Message",
			fields:   []int{},
			expected: "0000000000000000",
		},
		{
			name:     "Fields 3, 11, 41",
			fields:   []int{3, 11, 41},
			expected: "2020000000800000",
		},
		{
			name:   "Field 1 (Secondary Bitmap Indicator - Ignored in current impl)",
			fields: []int{1}, // Logic currently ignores < 1 || > 64, but 1 is handled as bit? Wait. Logic: fieldNum < 1 is ignored. 1 is valid?
			// byteIdx := (1 - 1) / 8 = 0. bitIdx := 7 - (0%8) = 7. 0x80.
			// Field 1 usually indicates secondary bitmap presence.
			// If we set field 1 manually, it sets the bit using the same logic.
			expected: "8000000000000000",
		},
		{
			name:   "Field 64 (Last bit of primary)",
			fields: []int{64},
			// byte 7. bitIdx = 7 - (63%8) = 0. 0x01.
			expected: "0000000000000001",
		},
		{
			name:   "High bit and Low bit (2 and 63)",
			fields: []int{2, 63},
			// Field 2: byte 0. bitIdx 6. 0x40.
			// Field 63: byte 7. bitIdx 1. 0x02.
			expected: "4000000000000002",
		}, {
			name:   "Field 65 (Secondary Bitmap)",
			fields: []int{65},
			// Expect 16 bytes.
			// Primary bitmap: Field 1 set (0x80...).
			// Secondary bitmap: Field 65 is first bit of second block (byte 8). 0x80.
			// hex: 8000000000000000 8000000000000000
			expected: "80000000000000008000000000000000",
		},
		{
			name:   "Field 128 (Last bit of Secondary)",
			fields: []int{128},
			// Primary: Field 1 set (0x80...).
			// Secondary: Field 128 is last bit (byte 15, bit 0). 0x01.
			expected: "80000000000000000000000000000001",
		},
		{
			name:   "Mixed Primary and Secondary (3, 65)",
			fields: []int{3, 65},
			// Primary: Field 1 (auto) + Field 3 (0x20). Byte 0: 0x80 | 0x20 = 0xA0.
			// Secondary: Field 65 (0x80 at byte 8).
			expected: "a0000000000000008000000000000000",
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMessage()
			for _, f := range tt.fields {
				m.Set(f, "dummy")
			}

			got, err := m.BitmapHex()
			if err != nil {
				t.Fatalf("BitmapHex() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("BitmapHex() = %v, want %v", got, tt.expected)
			}
		})
	}
}
