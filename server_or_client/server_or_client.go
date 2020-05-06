package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	//"../tlv_utils"
)

type ClientManager struct {
	clients map [*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data chan []byte
}

func startServerMode() {
	fmt.Println("Starting server...")
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println(err)
	}
	manager := ClientManager {
		clients: make(map[*Client]bool),
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
		manager.register <- client
		go manager.receive(client)
		go manager.send(client)
	}
}

func startClientMode() {
	fmt.Println("Starting client...")
	connection, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		fmt.Println(err)
	}
	client := &Client{socket: connection}
	go client.receive()
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')
		fmt.Printf("UNEDITED MESSAGE: |%s|\n", message)
		fmt.Printf("EDITED MESSAGE: |%s|\n",strings.TrimRight(message, "\n"))
		_, err := connection.Write([]byte(strings.TrimRight(message, "\n")))
		if err != nil {
			fmt.Println("ERROR WRITING TO CLIENT SIDE: ", err)
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
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED:", string(message[:length]))
			manager.broadcast <- message
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			_, err := client.socket.Write(message)

			if err != nil {
				fmt.Println("ERROR WRITING TO SOCKET:", err)
			}
		}
	}
}

func (client *Client) receive() {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED:", string(message))
		}
	}
}
