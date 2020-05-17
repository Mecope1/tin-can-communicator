package main

import (
	"flag"
	tccclient "github.com/Mecope1/tin-can-communicator/client"
	tccserver "github.com/Mecope1/tin-can-communicator/server"
	"strings"
)

func main() {
	// Server and client startup functions go here!
	flagMode := flag.String("mode", "server", "start in client or server mode")
	flagPort := flag.String("port", "8080", "set a port for your server to listen on")
	flagServerAddr := flag.String("address", "localhost:8080", "set an address for your client to dial into")
	flag.Parse()

	if strings.ToLower(*flagMode) == "server" {
		tccserver.StartServerMode(*flagPort)
	} else if strings.ToLower(*flagMode) == "client" {
		tccclient.StartClientMode(*flagServerAddr)
	}
}

