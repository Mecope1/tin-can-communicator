package client

import (
	"bufio"
	"fmt"
	bp "github.com/Mecope1/tin-can-communicator/binprot"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	socket      net.Conn
	data        chan []byte
	chatterName string
	roomID      uint
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
		fmt.Print("Enter a roomID (from 1 to 255): ")
		roomNumStr, err2 := myReader.ReadString('\n')

		if err1 != nil {
			fmt.Println(err1)
		}
		if err2 != nil {
			fmt.Println(err2)
		}

		if roomInt, err := strconv.Atoi(strings.TrimSpace(roomNumStr)); err == nil {
			if roomInt >= 1 && roomInt <= 255 && len(usrName) >= 3 {
				correctInput = true
				roomNum = uint(roomInt)
			}
		} else {
			fmt.Println("Incorrect username or room number.")
		}
	}

	// Store these values in the client object so that they can be passed around as needed.
	client.chatterName = usrName
	client.roomID = roomNum

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// BEGIN NEW UI WORK //

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewList()
	p.Rows = []string{"Hello World!"}

	termWidth, termHeight := ui.TerminalDimensions()
	p.Rows = append(p.Rows, "width: "+strconv.Itoa(termWidth)+" height: "+strconv.Itoa(termHeight))
	p.SetRect(0, 0, 100, 30)

	// Replace this with the joined channel's number.
	// num := 1

	p.Title = "Channel " + strconv.Itoa(int(roomNum))
	p.WrapText = true

	// These are included to move the cursor to the end of the list.
	// Otherwise the terminal will be 2 place behind the newest element, unless the user moves it themselves.
	p.ScrollDown()
	p.ScrollDown()

	q := widgets.NewParagraph()
	q.WrapText = true
	q.Title = client.chatterName
	q.SetRect(0, 0, 100, 4)

	grid := ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.6/2,
			ui.NewCol(1, p),
		),
		ui.NewRow(0.4/2,
			ui.NewCol(1, q),
		),
	)

	ui.Render(grid)

	// Triggers go routine that will listen for messages being broadcast from the server.
	go client.receive(p, grid)

	// Very basic authentication. Sends a record containing a specific magic byte to the server. If the server doesn't
	// receive that byte before other communications happen, then the server will close the connection.
	client.magicByteTest()

	for e := range ui.PollEvents() {
		if e.Type == ui.ResizeEvent {
			termWidth, termHeight := ui.TerminalDimensions()
			// This line redefines the size of the application, which is necessary when the terminal itself is resized.
			grid.SetRect(0, 0, termWidth, termHeight)
		} else if e.Type == ui.KeyboardEvent || e.Type == ui.MouseEvent {
			switch e.ID {
			case "<C-c>":
				fallthrough
			case "<Escape>":
				ui.Close()
				os.Exit(0)
			case "<PageUp>":
				p.ScrollPageUp()
			case "<PageDown>":
				p.ScrollPageDown()
			case "<Up>":
				fallthrough
			case "<MouseWheelUp>":
				p.ScrollUp()
			case "<Down>":
				fallthrough
			case "<MouseWheelDown>":
				p.ScrollDown()
			case "<Enter>":
				//p.Rows = append(p.Rows, q.Text)
				client.sendMessage(q)
				q.Text = ""
				//p.ScrollDown()
			case "<Backspace>":
				if len(q.Text) > 0 {
					chars := strings.Split(q.Text, "")

					q.Text = strings.Join(chars[0:len(chars)-1], "")
				}
			case "<Space>":
				q.Text += " "
			case "<Insert>":
			case "<Delete>":
			case "<Home>":
			case "<End>":
			case "<F1>":
			case "<F2>":
			case "<F3>":
			case "<F4>":
			case "<F5>":
			case "<F6>":
			case "<F7>":
			case "<F8>":
			case "<F9>":
			case "<F10>":
			case "<F11>":
			case "<F12>":
			case "<Left>":
			case "<Right>":
			case "<MouseRelease>":
			case "<MouseLeft>":
			case "<MouseMiddle>":
			case "<MouseRight>":
			default:
				q.Text += e.ID

			}

		}

		ui.Clear()
		ui.Render(grid)
	}

	// END NEW UI WORK //
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// Here the client will take input from the user and if it is well formed, it will write it to a channel that the
	// server receives from.
	//for {
	//	reader := bufio.NewReader(os.Stdin)
	//	msg, _ := reader.ReadString('\n')
	//	msg = strings.TrimRight(msg, "\n")
	//
	//	if len(msg) > 0 {
	//
	//		if len(msg) > 200 {
	//			fmt.Println("Message must be less than 200 characters. Message truncated.")
	//			msg = msg[:200]
	//		}
	//
	//		buf, err := bp.EncodeMsg(0x8A, msg, usrName, roomNum)
	//
	//		if err != nil {
	//			fmt.Println("Error during message construction: ", err.Error())
	//		} else {
	//			_, err := connection.Write(buf.Bytes())
	//			if err != nil {
	//				panic(err)
	//			}
	//		}
	//	}
	//}
}

func (client *Client) sendMessage(typingBox *widgets.Paragraph) {
	msg := typingBox.Text

	if len(msg) > 0 {

		if len(msg) > 200 {
			fmt.Println("Message must be less than 200 characters. Message truncated.")
			msg = msg[:200]
		}

		buf, err := bp.EncodeMsg(0x8A, msg, client.chatterName, client.roomID)

		if err != nil {
			fmt.Println("Error during message construction: ", err.Error())
		} else {
			_, err := client.socket.Write(buf.Bytes())
			if err != nil {
				panic(err)
			}
		}

		typingBox.Text = ""

	}
}

func (client *Client) magicByteTest() {

	// Special byte 0x8B is used rather than the usual 0x8A for standard messages. The server will reject the client
	// if any other bit is used!
	buf, err := bp.EncodeMsg(0x8B, "", client.chatterName, client.roomID)

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

func (client *Client) receive(chatBox *widgets.List, grid *ui.Grid) {
	//chatBox.Rows = append(chatBox.Rows, "LISTENING" )
	for {

		msgBytes := make([]byte, 256)
		length, err := client.socket.Read(msgBytes)
		//chatBox.Rows = append(chatBox.Rows, "NEW MESSAGE" )
		if err != nil {
			client.socket.Close()
			panic(err)
		}

		if length > 0 {
			newMessage, decodeErr := bp.DecodeMsg(msgBytes)
			if decodeErr != nil {
				fmt.Println("Error decoding message: ", decodeErr)
			} else {

				var textLine strings.Builder

				textLine.WriteString(string(newMessage.ChatterName))
				textLine.WriteString(": ")
				textLine.WriteString(string(newMessage.Payload))

				chatBox.Rows = append(chatBox.Rows, textLine.String())

			}
		}
		ui.Clear()
		ui.Render(grid)
	}
}
