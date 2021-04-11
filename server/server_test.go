package server

import (
	"fmt"
	"net"
	"testing"
	//tccclient "github.com/Mecope1/tin-can-communicator/client"
	bp "github.com/Mecope1/tin-can-communicator/binprot"
)


//func TestStartServerMode(t *testing.T) {
//	StartServerMode("12345")
//	// Output: Starting server. Listening on port 12345
//
//}

func TestCheckIncConnGood (t *testing.T) {
	var tcpClient net.Conn

	client := &Client{socket: tcpClient, data: make(chan []byte), chatterName: "N/A", roomID: 0}

	manager := &ClientManager {
		clients: make(map[*Client]bool),
		ChatRoomOccupants: make(map[uint][]string),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}

	usrName := "Mike"
	var roomID uint = 1

	buf, err := bp.EncodeMsg(0x8B, "", usrName, roomID)

	if err != nil {
		fmt.Println("Error during message construction: ", err.Error())
	} else {
		_, err := client.socket.Write(buf.Bytes())
		if err != nil {
			client.socket.Close()
			panic(err)
		}
	}

	isGoodClient := checkIncConn(client, manager)

	if isGoodClient == false {
		t.Errorf("checkIncConn failed to return true for a good client")
	}

}


func TestCheckIncConnBad (t *testing.T) {
	var tcpClient net.Conn
	client := &Client{socket: tcpClient, data: make(chan []byte), chatterName: "N/A", roomID: 0}
	manager := &ClientManager {
		clients: make(map[*Client]bool),
		ChatRoomOccupants: make(map[uint][]string),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}

	usrName := "Mike"
	var roomID uint = 1

	buf, err := bp.EncodeMsg(0x8C, "", usrName, roomID)

	if err != nil {
		fmt.Println("Error during message construction: ", err.Error())
	} else {
		_, err := client.socket.Write(buf.Bytes())
		if err != nil {
			client.socket.Close()
			panic(err)
		}
	}

	isGoodClient := checkIncConn(client, manager)

	if isGoodClient == true {
		t.Errorf("checkIncConn failed to return false for a bad client")
	}

}