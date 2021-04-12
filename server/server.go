package server

import (
	"encoding/hex"
	"fmt"
	bp "github.com/Mecope1/tin-can-communicator/binprot"
	"net"
)

type ClientManager struct {
	clients           map[*Client]bool
	ChatRoomOccupants map[uint][]string
	broadcast         chan []byte
	register          chan *Client
	unregister        chan *Client
}

type Client struct {
	socket      net.Conn
	data        chan []byte
	chatterName string
	roomID      uint
}

func StartServerMode(port string) {
	fmt.Println("Starting server. Listening on port", port)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
	}
	manager := ClientManager{
		clients:           make(map[*Client]bool),
		ChatRoomOccupants: make(map[uint][]string),
		broadcast:         make(chan []byte),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
	}

	go manager.start()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		client := &Client{socket: connection, data: make(chan []byte), chatterName: "N/A", roomID: 0}
		goodClient := checkIncConn(client, &manager)
		// authenticateUser()
		if goodClient {
			manager.register <- client
			go manager.receive(client)
			go manager.send(client)
		} else {
			client.socket.Close()
		}
	}
}

func checkIncConn(client *Client, man *ClientManager) bool {
	msgBytes := make([]byte, 256)
	retVal := false
	length, initReadErr := client.socket.Read(msgBytes)
	if length > 0 && initReadErr == nil {

		msg, err := bp.DecodeMsg(msgBytes)

		if err != nil {
			fmt.Println("Error decoding message: ", err.Error())
		} else if msg.Type == 139 {
			// Needs to be changed to check if the name already exists

			man.ChatRoomOccupants[msg.RoomID] = append(man.ChatRoomOccupants[msg.RoomID], string(msg.ChatterName))
			fmt.Printf("%s joined room: %d\n", string(msg.ChatterName), msg.RoomID)
			fmt.Printf("Room %d now has %d people in it\n", msg.RoomID, len(man.ChatRoomOccupants[msg.RoomID]))
			client.chatterName = string(msg.ChatterName)
			client.roomID = msg.RoomID
			retVal = true
		}
	}
	return retVal
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has been terminated!")
			}
		case message := <-manager.broadcast:
			record, err := bp.DecodeMsg(message)
			if err == nil {
				for connection := range manager.clients {
					if connection.roomID == record.RoomID {
						select {
						case connection.data <- message:

						default:
							close(connection.data)
							delete(manager.clients, connection)
						}
					}
				}
			} else {
				fmt.Println("Error processing client message: ", err.Error())
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
			fmt.Println("RECEIVED:", hex.Dump(message[:length]))
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
