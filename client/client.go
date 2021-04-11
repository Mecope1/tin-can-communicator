package client

import (
	"bufio"
	"fmt"
	bp "github.com/Mecope1/tin-can-communicator/binprot"
	"net"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	socket net.Conn
	data chan []byte
	chatterName string
	roomID uint
}


func StartClientMode(serverAddr string) {
	fmt.Println("Starting client. Dialing", serverAddr)
	connection, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	client := &Client{socket: connection}

	var usrName string
	var roomNum uint
	correctInput := false

	myReader := bufio.NewReader(os.Stdin)

	for !correctInput {
		fmt.Print("Enter a username (3 or more chars): ")
		usrNameInpt, err1 := myReader.ReadString('\n')
		usrName = strings.TrimSpace(usrNameInpt)
		fmt.Print("Enter a roomID (from 1 to 255): " )
		roomNumStr, err2 := myReader.ReadString('\n')

		if err1 != nil {
			fmt.Println(err1)
		}
		if err2 != nil {
			fmt.Println(err2)
		}

		if roomInt, err := strconv.Atoi(strings.TrimSpace(roomNumStr)); err == nil {
			if roomInt >= 1 && roomInt <= 255 && len(usrName) >= 3  {
				correctInput = true
				roomNum = uint(roomInt)
			}
		} else {
			fmt.Println("Incorrect username or room number.")
		}
	}

	// Triggers go routine that will listen for messages being broadcast from the server.
	go client.receive()

	// Very basic authentication. Sends a record containing a specific magic byte to the server. If the server doesn't
	// receive that byte before other communications happen, then the server will close the connection.
	magicByteTest(client, usrName, roomNum)

	// Here the client will take input from the user and if it is well formed, it will write it to a channel that the
	// server receives from.
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimRight(msg, "\n")

		if len(msg) > 0 {

			if len(msg) > 200 {
				fmt.Println("Message must be less than 200 characters. Message truncated.")
				msg = msg[:200]
			}

			buf, err:= bp.EncodeMsg(0x8A, msg, usrName, roomNum)

			if err != nil {
				fmt.Println("Error during message construction: ", err.Error())
			} else {
				_, err := connection.Write(buf.Bytes())
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func magicByteTest(client *Client, usrName string, roomID uint) {

	// Special byte 0x8B is used rather than the usual 0x8A for standard messages. The server will reject the client
	// if any other bit is used!
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
}

func (client *Client) receive() {

	for {
		msgBytes := make([]byte, 256)
		length, err := client.socket.Read(msgBytes)
		if err != nil {
			client.socket.Close()
			panic(err)
		}

		if length > 0 {
			msg, decodeErr := bp.DecodeMsg(msgBytes)
			if decodeErr != nil {
				fmt.Println("Error decoding message: ", decodeErr)
			} else if string(msg.ChatterName) != client.chatterName {
				fmt.Printf("In Room#%d |%s| Said: %s\n", msg.RoomID, string(msg.ChatterName), string(msg.Payload))
			}
		}
	}
}
