package main

import (
	"encoding/json"
	"os"
	"proj2/server"
	"strconv"
)

func main() {
	args := os.Args[1:]

	// Determine mode and number of consumers
	mode := "s"
	consumers := 0
	if len(args) > 0 {
		mode = "p"
		consumers, _ = strconv.Atoi(args[0])
	}

	// Create server configuration
	config := server.Config{
		Encoder:        json.NewEncoder(os.Stdout),
		Decoder:        json.NewDecoder(os.Stdin),
		Mode:           mode,
		ConsumersCount: consumers,
	}

	// Run the server
	server.Run(config)
}
