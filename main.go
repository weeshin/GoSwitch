package main

import (
	"GoSwitch/pkg/client"
	"GoSwitch/pkg/config"
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

	// 2. Load ISO Spec (Field definitions)
	isoSpec, err := iso8583.LoadSpecFromFile("default_spec.yaml")
	if err != nil {
		log.Fatalf("Error loading spec.yaml: %v", err)
	}

	// 3. Start Server using config values
	addr := fmt.Sprintf("%s:%d", appCfg.Server.IP, appCfg.Server.Port)
	srv := server.New(addr, isoSpec)

	// 4. Start Clients for all enabled channels
	for _, ch := range appCfg.Channels {
		if !ch.Enabled {
			continue
		}
		// Launch client coroutine
		go startClient(ch, isoSpec)
	}

	log.Printf("Starting GoSwitch on %s", addr)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}

func startClient(ch config.ChannelConfig, spec iso8583.Spec) {
	c := client.New(ch, spec)
	// Connect blocks until successful
	if err := c.Connect(); err == nil {
		// Send a Test Network Message (0800) once connected
		echo := iso8583.NewMessage()
		echo.MTI = "0800"
		echo.Set(70, "301")
		if err := c.Send(echo); err != nil {
			log.Printf("Failed to send echo to %s: %v", c.Config.Name, err)
		}
	}
}
