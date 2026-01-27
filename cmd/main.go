package main

import (
	"GoSwitch/pkg/config"
	"GoSwitch/pkg/field"
	"GoSwitch/pkg/iso8583"
	"GoSwitch/pkg/server"
	"fmt"
	"log"
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
	}
	spec.Fields = map[int]iso8583.FieldSpec{
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
		39: {
			Length:      2,
			Description: "Response Code",
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
			Encoder:     &field.FALLLChar{},
		},
		55: {
			Length:      255,
			Description: "ICC Data – EMV Having Multiple Tags",
			Encoder:     &field.FALLLChar{},
		},
		56: {
			Length:      999,
			Description: "Private Field",
			Encoder:     &field.FALLLChar{},
		},
		57: {
			Length:      999,
			Description: "Private Field (NATIONAL)",
			Encoder:     &field.FALLLChar{},
		},
		58: {
			Length:      999,
			Description: "Private Field (NATIONAL)",
			Encoder:     &field.FALLLChar{},
		},
		59: {
			Length:      999,
			Description: "Private Field (NATIONAL)",
			Encoder:     &field.FALLLChar{},
		},
		60: {
			Length:      999,
			Description: "Private Field",
			Encoder:     &field.FALLLChar{},
		},
		61: {
			Length:      999,
			Description: "Private Field",
			Encoder:     &field.FALLLChar{},
		},
		62: {
			Length:      999,
			Description: "Private Field",
			Encoder:     &field.FALLLChar{},
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
	}

	addr := fmt.Sprintf("%s:%d", appCfg.Server.IP, appCfg.Server.Port)
	app := server.NewEngine(addr, spec, &server.NACChannel{})

	// 3. Define your Logic (The app.Request handler)
	app.Request(func(c *iso8583.Context) {
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
func handleEcho(c *iso8583.Context) {
	resp := iso8583.NewMessage()
	resp.MTI = "0810"
	// Copy the STAN (Field 11) from request to response
	if stan, ok := c.Request.Fields[11]; ok {
		resp.Set(11, string(stan.Value))
	}
	resp.Set(39, "00") // Action Code: Approved
	c.Respond(resp)
}

func handlePurchase(c *iso8583.Context) {
	// Implement your database or authorization logic here
	fmt.Println("Processing Purchase...")
	resp := iso8583.NewMessage()
	resp.MTI = "0210"
	if stan, ok := c.Request.Fields[11]; ok {
		resp.Set(11, string(stan.Value))
	}
	if tid, ok := c.Request.Fields[41]; ok {
		resp.Set(41, string(tid.Value))
	}
	resp.Set(39, "00") // Action Code: Approved
	c.Respond(resp)
}
