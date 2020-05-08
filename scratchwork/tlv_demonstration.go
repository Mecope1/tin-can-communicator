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

	record := &tlv_utils.Record {
		Type: 0x8A,
		Version: 16, // 0x45,
		RoomID:  3, // 0x1,
		ChatterID: 255, // 0xff,
		ChatterName: []byte("Mike"),
		Payload: []byte("hello, go! DID YOU SEE THAT?!?!?!?!?!?!?!?!?"),
	}

	err	:= tlvWriter.Write(record)

	if err != nil {
		fmt.Println ("ERROR WRITING: ", err)
	}

	fmt.Println(hex.Dump(buf.Bytes()))

	reader := bytes.NewReader(buf.Bytes())
	tlvReader := tlv_utils.NewReader(reader, codec)

	next, readErr := tlvReader.Next()
	if readErr != nil {
		fmt.Println("reading error: ", readErr)
	}

	fmt.Printf("type: |%d|\n", next.Type)
	fmt.Printf("version: |%d|\n", next.Version)
	fmt.Printf("roomID: |%d|\n", next.RoomID)
	fmt.Printf("chatterID: |%d|\n", next.ChatterID)
	fmt.Printf("chatterName: |%s|\n", string(next.ChatterName))

	fmt.Printf("payload: |%s|\n", string(next.Payload))
}