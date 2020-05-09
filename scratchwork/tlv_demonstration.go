package main

import (
	"../tlv_utils"
	"bytes"
	"encoding/hex"
	"fmt"
)

func main() {

	buf := new(bytes.Buffer)
	codec := &tlv_utils.Codec{
		TypeBytes: tlv_utils.OneByte,
		VersionBytes: tlv_utils.OneByte,
		ChatterIDBytes: tlv_utils.OneByte,
		ChatterNameBytes: tlv_utils.TwoBytes,
		RoomIDBytes: tlv_utils.OneByte,
		PayloadBytes: tlv_utils.FourBytes}

	tlvWriter := tlv_utils.NewWriter(buf, codec)

	/*
	The record contains information about the sender, the type of message, and other information including a
	message payload
	*/

	record := &tlv_utils.Record {
		Type: 0x8A,
		Version: 16,
		RoomID:  3,
		ChatterID: 255,
		ChatterName: []byte("Mike"),
		Payload: []byte("hello, go! DID YOU SEE THAT?!?!?!?!?!?!?!?!?"),
	}

	/*
	Here we write a record that follows the pattern of the codec into a buffer
	that was specified when the writer was made.
	*/
	err	:= tlvWriter.Write(record)
	if err != nil {
		fmt.Println ("ERROR WRITING: ", err)
	}

	// This step isn't necessary, but shows the state of the data once it has been encoded.
	fmt.Println(hex.Dump(buf.Bytes()))



	// Now is when the record is actually decoded. The reader also requires a codec to decode the []byte.

	reader := bytes.NewReader(buf.Bytes())
	tlvReader := tlv_utils.NewReader(reader, codec)

	/*
	Here we write out the contents of the []bytes into a next, which is a "record".*/
	next, readErr := tlvReader.Next()
	if readErr != nil {
		fmt.Println("reading error: ", readErr)
	}

	// This just shows us how the data can be accessed, and that it made its journey unharmed.
	fmt.Printf("type: |%d|\n", next.Type)
	fmt.Printf("version: |%d|\n", next.Version)
	fmt.Printf("roomID: |%d|\n", next.RoomID)
	fmt.Printf("chatterID: |%d|\n", next.ChatterID)
	fmt.Printf("chatterName: |%s|\n", string(next.ChatterName))

	fmt.Printf("payload: |%s|\n", string(next.Payload))
}