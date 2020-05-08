package main

import (
	"../tlv_utils"
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var constCodec = &tlv_utils.Codec{
		TypeBytes: tlv_utils.OneByte,
		VersionBytes: tlv_utils.OneByte,
		ChatterIDBytes: tlv_utils.OneByte,
		ChatterNameBytes: tlv_utils.TwoBytes,
		RoomIDBytes: tlv_utils.OneByte,
		PayloadBytes: tlv_utils.FourBytes,
}

type ClientManager struct {
	clients map [*Client]bool
	ChatRoomOccupants map [uint][]string
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data chan []byte
	//ChatterID uint
	//chatterName string
}

func startServerMode() {
	fmt.Println("Starting server...")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	manager := ClientManager {
		clients: make(map[*Client]bool),
		ChatRoomOccupants: make(map[uint][]string),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client)}

	go manager.start()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		client := &Client{socket:connection, data: make(chan []byte)}
		goodClient := checkIncConn(client, &manager)
		// authenticateUser()
		if goodClient == true {
			manager.register <- client
			go manager.receive(client)
			go manager.send(client)
		} else if goodClient == false {
			client.socket.Close()
		}
	}
}

func checkIncConn (client *Client, man *ClientManager) bool {
	msg := make([]byte, 256)
	retVal := false
	length, initReadErr := client.socket.Read(msg)
	if length > 0 && initReadErr == nil {
		reader := bytes.NewReader(msg)
		tlvReader := tlv_utils.NewReader(reader, constCodec)
		next, readingError := tlvReader.Next()
		if readingError != nil {
			fmt.Println("Error decoding message: ", readingError)
		} else if next.Type == 139 {
			// Needs to be changed to check if the name already exists, else it will spam both people twice at least.
				man.ChatRoomOccupants[next.RoomID] = append(man.ChatRoomOccupants[next.RoomID], string(next.ChatterName))
				fmt.Printf("%s joined room: %d\n", string(next.ChatterName), next.RoomID)
				fmt.Printf("Room %d now has %d people in it\n", next.RoomID, len(man.ChatRoomOccupants[next.RoomID]))
			retVal = true
		}
	}
	return retVal
}


//  MAKE INTO SERVER-SIDE AUTH FN
//func authenticateUser() {
//	client := &Client{socket:connection, data: make(chan []byte)}
//
//
//	codec := &tlv_utils.Codec{
//		TypeBytes: tlv_utils.OneByte,
//		VersionBytes: tlv_utils.OneByte,
//		ChatterIDBytes: tlv_utils.OneByte,
//		ChatterNameBytes: tlv_utils.TwoBytes,
//		RoomIDBytes: tlv_utils.OneByte,
//		PayloadBytes: tlv_utils.FourBytes,
//	}
//
//	msg := make([]byte, 256)
//	length, initReadErr := connection.Read(msg)
//	if length > 0 && initReadErr == nil {
//		reader := bytes.NewReader(msg)
//		tlvReader := tlv_utils.NewReader(reader,  constCodec)
//		next, readingError := tlvReader.Next()
//		if readingError != nil {
//			fmt.Println("Error decoding message: ", readingError)
//		} else if readingError == nil && next.Type == 138{
//			manager.register <- client
//			go manager.receive(client)
//			go manager.send(client)
//		}
//
//	} else {
//
//		fmt.Println("Incorrect magic byte given by: ", connection.RemoteAddr())
//
//		buf := new(bytes.Buffer)
//		codec := &tlv_utils.Codec{
//			TypeBytes: tlv_utils.OneByte,
//			VersionBytes: tlv_utils.OneByte,
//			ChatterIDBytes: tlv_utils.OneByte,
//			ChatterNameBytes: tlv_utils.TwoBytes,
//			RoomIDBytes: tlv_utils.OneByte,
//			PayloadBytes: tlv_utils.FourBytes}
//
//		tlvWriter := tlv_utils.NewWriter(buf,  constCodec)
//
//		errorMsg := &tlv_utils.Record{
//			Type:        0x8A,
//			Version:     1,
//			RoomID:      0,
//			ChatterID:   0,
//			ChatterName: []byte("Server"),
//			Payload:     []byte("Incorrect Magic Byte. Please ensure you are connecting to the correct server!!!"),
//		}
//
//		tlvWrErr := tlvWriter.Write(errorMsg)
//
//		if tlvWrErr != nil {
//			fmt.Println("Error during message construction: ", tlvWrErr.Error())
//			continue
//		} else {
//			_, connWrErr := connection.Write(buf.Bytes())
//			if connWrErr != nil {
//				fmt.Println("Error writing to client channel: ", connWrErr)
//			}
//		}
//		manager.unregister <- client
//	}
//}


func startClientMode() {
	fmt.Println("Starting client...")
	connection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
		// fmt.Println(err)
	}
	client := &Client{socket: connection}

	//buf := new(bytes.Buffer)
	//tlvWriter := tlv_utils.NewWriter(buf,  constCodec)


	// temp_const
	//codec := &tlv_utils.Codec{
	//	TypeBytes: tlv_utils.OneByte,
	//	VersionBytes: tlv_utils.OneByte,
	//	ChatterIDBytes: tlv_utils.OneByte,
	//	ChatterNameBytes: tlv_utils.TwoBytes,
	//	RoomIDBytes: tlv_utils.OneByte,
	//	PayloadBytes: tlv_utils.FourBytes}



	var usrName string
	var roomNum uint
	var correctInput bool = false

	myReader := bufio.NewReader(os.Stdin)

	for !correctInput {
		fmt.Print("Enter a username (3 or more chars): ")
		usrNameInpt, err1 := myReader.ReadString('\n')
		usrName = strings.TrimSpace(usrNameInpt)
		fmt.Print("Enter a roomID (unsigned ints only!): " )
		roomNumStr, err2 := myReader.ReadString('\n')

		if err1 != nil {
			fmt.Println(err1)
		}

		if err2 != nil {
			fmt.Println(err2)
		}

		if roomInt, err := strconv.Atoi(strings.TrimSpace(roomNumStr)); err == nil {
			if roomInt > 0 && len(usrName) >= 3  {
				correctInput = true
				roomNum = uint(roomInt)
			}
		} else {
			fmt.Println("Incorrect username or room number.")
		}
	}

	go client.receive()

	// Auth should happen here
	magicByteTest(client, &usrName, &roomNum)

	// Attempt to use ANSI escape chars
	//var entryField *bytes.Buffer = new(bytes.Buffer)
	//var usrOutput *bufio.Writer = bufio.NewWriter(os.Stdout)

	for {
		buf := new(bytes.Buffer)
		tlvWriter := tlv_utils.NewWriter(buf, constCodec)
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')
		//fmt.Println("THING GOES HERE")
		msg = strings.TrimRight(msg, "\n")

		if len(msg) > 0 {
			record := &tlv_utils.Record{
				Type:        0x8A,
				Version:     1,
				RoomID:      roomNum,
				ChatterID:   255,
				ChatterName: []byte(usrName),
				Payload:     []byte(msg),
			}

			err := tlvWriter.Write(record)

			if err != nil {
				fmt.Println("Error during message construction: ", err.Error())
				continue
			} else {
				_, err := connection.Write(buf.Bytes())
				if err != nil {
					panic(err)
				}
			}

			// Attempt to use ANSI escape chars
			//usrOutput.WriteString("\033[2J")
			//fmt.Fprintf(entryField, "\033[1;1H")

			// Pre-TLV usage
			//_, err := connection.Write([]byte(strings.TrimRight(message, "\n")))
			//if err != nil {
			//	panic(err)
			//}
		}
	}
}

// TURN INTO CLIENTSIDE AN AUTH FN
func magicByteTest(client *Client, usrName *string, roomID *uint) {
	buf := new(bytes.Buffer)
	tlvWriter := tlv_utils.NewWriter(buf, constCodec)
	rec := &tlv_utils.Record {

		Type: 0x8B,
		Version: 1,
		RoomID: *roomID,
		ChatterID: 0,
		ChatterName: []byte(*usrName),
		Payload: []byte(""),
	}

	err := tlvWriter.Write(rec)
	if err != nil {
		fmt.Println("Error during message construction: ", err.Error())
	} else {
		_, err := client.socket.Write(buf.Bytes())
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	flagMode := flag.String("mode", "server", "start in client or server mode")
	flag.Parse()
	if strings.ToLower(*flagMode) == "server" {
		startServerMode()
	} else {
		startClientMode()
	}
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <- manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has been terminated!")
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:

				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 256)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}

		if length > 0 {
			fmt.Println("RECEIVED:", string(hex.Dump(message[:length])))
			manager.broadcast <- message
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case msg, ok := <-client.data:
			if !ok {
				return
			}

			_, err := client.socket.Write(msg)

			if err != nil {
				fmt.Println("ERROR WRITING TO SOCKET:", err)
			}
		}
	}
}

func (client *Client) receive() {

	//buf := new(bytes.Buffer)
	//temp_const
	// constCodec := &tlv_utils.Codec{
	//	TypeBytes: tlv_utils.OneByte,
	//	VersionBytes: tlv_utils.OneByte,
	//	ChatterIDBytes: tlv_utils.OneByte,
	//	ChatterNameBytes: tlv_utils.TwoBytes,
	//	RoomIDBytes: tlv_utils.OneByte,
	//	PayloadBytes: tlv_utils.FourBytes}

	for {
		msg := make([]byte, 256)
		length, err := client.socket.Read(msg)
		if err != nil {
			client.socket.Close()
			break
		}

		if length > 0 {
			reader := bytes.NewReader(msg)
			tlvReader := tlv_utils.NewReader(reader, constCodec)
			next, readingError := tlvReader.Next()
			if readingError != nil {
				fmt.Println("Error decoding message: ", readingError)
			} else  {
				fmt.Printf("In Room#%d |%s| Said: %s\n", next.RoomID, string(next.ChatterName), string(next.Payload))
			}
		}
	}
}
