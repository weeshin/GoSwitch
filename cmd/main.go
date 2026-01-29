package main

import (
	"GoSwitch/pkg/config"
	"GoSwitch/pkg/field"
	"GoSwitch/pkg/iso8583"
	"GoSwitch/pkg/server"
	"fmt"
	"log"
	"time"
)

func main() {
	// 1. Load Application Config (Ports, etc.)
	appCfg, err := config.LoadAppConfig("app.yaml")
	if err != nil {
		log.Fatalf("Error loading app.yaml: %v", err)
	}

	// spec, _ := iso8583.LoadSpecFromFile("spec.yaml")

	spec := &iso8583.Spec{
		MTIEncoder:    &field.FBNumeric{},
		BitmapEncoder: &field.FBBitmap{},
		Fields: map[int]iso8583.FieldSpec{
			2: {
				Length:      16,
				Description: "Primary Account Number",
				Encoder:     &field.FBLLNumeric{},
			},
			3: {
				Length:      6,
				Description: "Processing Code",
				Encoder:     &field.FBNumeric{},
			},
			4: {
				Length:      12,
				Description: "Amount, Transaction",
				Encoder:     &field.FBNumeric{},
			},
			5: {
				Length:      12,
				Description: "Amount, Settlement",
				Encoder:     &field.FBNumeric{},
			},
			6: {
				Length:      12,
				Description: "Amount, Cardholder Billing",
				Encoder:     &field.FBNumeric{},
			},
			11: {
				Length:      6,
				Description: "Systems Trace Audit Number",
				Encoder:     &field.FBNumeric{},
			},
			12: {
				Length:      6,
				Description: "Time, Local Transaction",
				Encoder:     &field.FBNumeric{},
			},
			13: {
				Length:      4,
				Description: "Date, Local Transaction",
				Encoder:     &field.FBNumeric{},
			},
			14: {
				Length:      4,
				Description: "Date, Expiration",
				Encoder:     &field.FBNumeric{},
			},
			15: {
				Length:      4,
				Description: "Date, Settlement",
				Encoder:     &field.FBNumeric{},
			},
			16: {
				Length:      4,
				Description: "Date, Conversion",
				Encoder:     &field.FBNumeric{},
			},
			17: {
				Length:      4,
				Description: "Date, Capture",
				Encoder:     &field.FBNumeric{},
			},
			18: {
				Length:      4,
				Description: "Merchant Type",
				Encoder:     &field.FBNumeric{},
			},
			19: {
				Length:      3,
				Description: "Acquiring Institution Country Code",
				Encoder:     &field.FBNumeric{},
			},
			20: {
				Length:      3,
				Description: "PAN Country Code",
				Encoder:     &field.FBNumeric{},
			},
			21: {
				Length:      3,
				Description: "Forwarding Institution Country Code",
				Encoder:     &field.FBNumeric{},
			},
			22: {
				Length:      3,
				Description: "Point of Service Entry Mode",
				Encoder:     &field.FBNumeric{},
			},
			23: {
				Length:      3,
				Description: "Card Sequence Number",
				Encoder:     &field.FBNumeric{},
			},
			24: {
				Length:      3,
				Description: "Network International Identifier",
				Encoder:     &field.FBNumeric{},
			},
			25: {
				Length:      2,
				Description: "Point of Service Condition Code",
				Encoder:     &field.FBNumeric{},
			},
			26: {
				Length:      2,
				Description: "Point of Service PIN Capture Code",
				Encoder:     &field.FBNumeric{},
			},
			27: {
				Length:      1,
				Description: "Authorization Identification Response Length",
				Encoder:     &field.FBNumeric{},
			},
			28: {
				Length:      9,
				Description: "Amount, Transaction Fee",
				Encoder:     &field.FBNumeric{},
			},
			29: {
				Length:      9,
				Description: "Amount, Settlement Fee",
				Encoder:     &field.FBNumeric{},
			},
			30: {
				Length:      9,
				Description: "Amount, Transaction Processing Fee",
				Encoder:     &field.FBNumeric{},
			},
			31: {
				Length:      9,
				Description: "Amount, Settlement Processing Fee",
				Encoder:     &field.FBNumeric{},
			},
			32: {
				Length:      11,
				Description: "Acquiring Institution Identification Code",
				Encoder:     &field.FBLLNumeric{},
			},
			33: {
				Length:      11,
				Description: "Forwarding Institution Identification Code",
				Encoder:     &field.FBLLNumeric{},
			},
			34: {
				Length:      28,
				Description: "Primary Account Number, Extended",
				Encoder:     &field.FBLLChar{},
			},
			35: {
				Length:      37,
				Description: "Track 2 Data",
				Encoder:     &field.FBLLNumeric{},
			},
			36: {
				Length:      104,
				Description: "Track 3 Data",
				Encoder:     &field.FBLLLChar{},
			},
			37: {
				Length:      12,
				Description: "Retrieval Reference Number",
				Encoder:     &field.FChar{},
			},
			38: {
				Length:      6,
				Description: "Authorization Identification Response",
				Encoder:     &field.FChar{},
			},
			39: {
				Length:      2,
				Description: "Response Code",
				Encoder:     &field.FChar{},
			},
			40: {
				Length:      3,
				Description: "Service Restriction Code",
				Encoder:     &field.FChar{},
			},
			41: {
				Length:      8,
				Description: "Card Acceptor Terminal Identification",
				Encoder:     &field.FChar{},
			},
			42: {
				Length:      15,
				Description: "Card Acceptor Identification Code",
				Encoder:     &field.FChar{},
			},
			43: {
				Length:      40,
				Description: "Card Acceptor Name/Location",
				Encoder:     &field.FChar{},
			},
			44: {
				Length:      25,
				Description: "Additional Response Data",
				Encoder:     &field.FBLLChar{},
			},
			45: {
				Length:      76,
				Description: "Track 1 Data",
				Encoder:     &field.FBLLChar{},
			},
			46: {
				Length:      999,
				Description: "Additional Data - ISO",
				Encoder:     &field.FBLLLChar{},
			},
			47: {
				Length:      999,
				Description: "Additional Data - National",
				Encoder:     &field.FBLLLChar{},
			},
			48: {
				Length:      999,
				Description: "Additional Data - Private",
				Encoder:     &field.FBLLLChar{},
			},
			49: {
				Length:      3,
				Description: "Currency Code, Transaction",
				Encoder:     &field.FChar{},
			},
			50: {
				Length:      3,
				Description: "Currency Code, Settlement",
				Encoder:     &field.FChar{},
			},
			51: {
				Length:      3,
				Description: "Currency Code, Cardholder Billing",
				Encoder:     &field.FChar{},
			},
			52: {
				Length:      8,
				Description: "Personal Identification Number (PIN) Data",
				Encoder:     &field.FBBinary{},
			},
			53: {
				Length:      16,
				Description: "Security Related Control Information",
				Encoder:     &field.FBNumeric{},
			},
			54: {
				Length:      120,
				Description: "Additional Amounts",
				Encoder:     &field.FBLLLChar{},
			},
			55: {
				Length:      255,
				Description: "ICC Data – EMV Having Multiple Tags",
				Encoder:     &field.FBLLLChar{},
			},
			56: {
				Length:      999,
				Description: "Private Field",
				Encoder:     &field.FBLLLChar{},
			},
			57: {
				Length:      999,
				Description: "Private Field (NATIONAL)",
				Encoder:     &field.FBLLLChar{},
			},
			58: {
				Length:      999,
				Description: "Private Field (NATIONAL)",
				Encoder:     &field.FBLLLChar{},
			},
			59: {
				Length:      999,
				Description: "Private Field (NATIONAL)",
				Encoder:     &field.FBLLLChar{},
			},
			60: {
				Length:      999,
				Description: "Private Field",
				Encoder:     &field.FBLLLChar{},
			},
			61: {
				Length:      999,
				Description: "Private Field",
				Encoder:     &field.FBLLLChar{},
			},
			62: {
				Length:      999,
				Description: "Private Field",
				Encoder:     &field.FBLLLChar{},
			},
			63: {
				Length:      999,
				Description: "Private Field",
				Encoder:     &field.FBLLLChar{},
			},
			64: {
				Length:      16,
				Description: "Message Authentication Code (MAC)",
				Encoder:     &field.FBBinary{},
			},
		},
	}
	addr := fmt.Sprintf("%s:%d", appCfg.Server.IP, appCfg.Server.Port)
	// Example TPDU for a specific bank: 60 00 01 00 00
	bankTPDU := []byte{0x60, 0x00, 0x01, 0x00, 0x00}

	// Create NAC Channel with specific TPDU and Handler
	nacChannel := &server.NACChannel{
		BaseChannel: &server.BaseChannel{
			Spec:    spec,
			Header:  bankTPDU,
			Handler: &server.NACHeader{},
		},
	}

	app := server.NewEngine(addr, spec, nacChannel)

	// 3. Define your Logic (The app.Request handler)
	app.Request(func(c *server.Context) {
		// Define the combination key: MTI + Field 3
		mti := c.Request.MTI
		procCode := ""
		if f3, ok := c.Request.Fields[3]; ok {
			procCode = string(f3.Value)
		}

		key := fmt.Sprintf("%s_%s", mti, procCode)
		fmt.Printf("--> Handling Transaction: [%s]\n", key)

		switch key {
		case "0800_": // Network Echo
			handleEcho(c)
		case "0200_000000": // Purchase
			handlePurchase(c)
		default:
			fmt.Printf("No specific handler for %s\n", key)
		}
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// Logic for Echo
func handleEcho(c *server.Context) {
	resp := iso8583.NewMessage()
	resp.MTI = "0810"
	// Copy the STAN (Field 11) from request to response
	if stan, ok := c.Request.Fields[11]; ok {
		resp.Set(11, string(stan.Value))
	}
	resp.Set(39, "00") // Action Code: Approved
	c.Channel.Send(c.Conn, resp)
}

func handlePurchase(c *server.Context) {
	// Implement your database or authorization logic here
	fmt.Println("Processing Purchase...")

	c.Request.ResponseMTI()
	resp := c.Request
	resp.Set(39, "00")
	time.Sleep(1 * time.Second)
	c.Channel.Send(c.Conn, resp)
}
