package client

import "testing"

func TestStartClientMode(t *testing.T) {

	StartClientMode("127.0.0.1:8080")
	// Output: Starting client. Dialing 127.0.0.1:8080


}