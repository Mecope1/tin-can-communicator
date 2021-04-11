package main

import (
	"flag"
	tccclient "github.com/Mecope1/tin-can-communicator/client"
	tccserver "github.com/Mecope1/tin-can-communicator/server"
	"strings"
)


func main() {
	// Server and client startup functions go here!
	// The two options for this flag are server or client, depending on what function is desired
	flagMode := flag.String("mode", "client", "start in client or server mode")

	// Here, either the port for the server to listen on, or the address and port to dial into are set.

	// Setting the wrong flag for your desired mode will simply use the default values for the correct flag if no value
	// was given for it.
	flagPort := flag.String("port", "8080", "set a port for your server to listen on")
	flagServerAddr := flag.String("address", "localhost:8080", "set an address for your client to dial into")
	flag.Parse()

	if strings.ToLower(*flagMode) == "server" {
		tccserver.StartServerMode(*flagPort)
	} else if strings.ToLower(*flagMode) == "client" {
		tccclient.StartClientMode(*flagServerAddr)
	}
}

