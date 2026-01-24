package main

import (
	"GoSwitch/pkg/iso8583"
	"fmt"
)

func main() {
	msg := iso8583.NewMessage()
	msg.Set(3, "000000")   // Processing Code
	msg.Set(11, "123456")  // STAN
	msg.Set(41, "TERM001") // Terminal ID

	hexStr, _ := msg.BitmapHex()
	fmt.Printf("Bitmap: %s\n", hexStr)
	// Output should show bits 3, 11, and 41 as '1'
}
