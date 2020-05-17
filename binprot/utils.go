package binprot

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type ByteSize int

const (
	OneByte ByteSize = 1
	TwoBytes ByteSize = 2
	FourBytes ByteSize = 4
	EightBytes ByteSize = 8
)

var RecordCodec = &Codec{
	TypeBytes: OneByte,
	VersionBytes: OneByte,
	ChatterIDBytes: OneByte,
	ChatterNameBytes: TwoBytes,
	RoomIDBytes: OneByte,
	PayloadBytes: FourBytes,
}

type Record struct {
	Type uint
	Version uint
	RoomID uint
	ChatterID uint
	ChatterName []byte
	Payload []byte
}

type Codec struct {
	TypeBytes        ByteSize // Number of bytes used for the type field
	VersionBytes     ByteSize // Number of bytes used for the version field
	RoomIDBytes      ByteSize // Number of bytes used for the chat room ID field
	ChatterIDBytes   ByteSize // Number of bytes used for the host ID field
	ChatterNameBytes ByteSize // Number of bytes used for the host name field
	PayloadBytes     ByteSize // Number of bytes used for the length field
}

// WRITER SECTION
type TLVWriter struct {
	writer io.Writer
	codec *Codec
}


func newWriter( w io.Writer, codec *Codec) *TLVWriter {
	newWriter := TLVWriter{w, codec}
	return &newWriter
}


func (w *TLVWriter) write(rec *Record) error {
	err := writeUint(w.writer, w.codec.TypeBytes, rec.Type)
	if err != nil {
		return err
	}


	err = writeUint(w.writer, w.codec.VersionBytes, rec.Version)
	if err != nil {
		return err
	}

	err = writeUint(w.writer, w.codec.RoomIDBytes, rec.RoomID)
	if err != nil {
		return err
	}

	err = writeUint(w.writer, w.codec.ChatterIDBytes, rec.ChatterID)
	if err != nil {
		return err
	}

	uchatterNameLen := uint(len(rec.ChatterName))
	err = writeUint(w.writer, w.codec.ChatterNameBytes, uchatterNameLen)
	if err != nil {
		return err
	}

	_, err = w.writer.Write(rec.ChatterName)
	if err != nil {
		return err
	}

	ulen := uint(len(rec.Payload))
	err = writeUint(w.writer, w.codec.PayloadBytes, ulen)
	if err != nil {
		return err
	}

	_, err = w.writer.Write(rec.Payload)
	return err
}


func writeUint(w io.Writer, b ByteSize, i uint) error {
	var num interface{}
	switch b {
	case OneByte:
		num=uint8(i)
	case TwoBytes:
		num=uint16(i)
	case FourBytes:
		num=uint32(i)
	case EightBytes:
		num=uint64(i)
	}

	return binary.Write(w, binary.BigEndian, num)
}


// READER SECTION
type Reader struct {
	codec *Codec
	reader io.Reader
}


func newReader(reader io.Reader, codec *Codec) *Reader {
	newReader := Reader{codec, reader}
	return &newReader
}


func(r *Reader) next() (*Record, error) {

	// Read in the bytes that make up the type field
	typeBytes := make([]byte, r.codec.TypeBytes)
	_, err := r.reader.Read(typeBytes)
	if err != nil {
		return nil, err
	}
	recType := readUint(typeBytes, r.codec.TypeBytes)

	// Read in the bytes that make up the version field
	versionBytes := make([]byte, r.codec.VersionBytes)
	_, err = r.reader.Read(versionBytes)
	if err != nil && err != io.EOF {
		return nil, err
	}
	recVersion := readUint(versionBytes, r.codec.VersionBytes)

	// Read in the bytes that make up the RoomID field
	roomIDBytes := make([]byte, r.codec.RoomIDBytes)
	_, err = r.reader.Read(roomIDBytes)
	if err != nil && err != io.EOF {
		return nil, err
	}
	recRoomID := readUint(roomIDBytes, r.codec.RoomIDBytes)

	// Read in the bytes that make up the ChatterID field
	chatterIDBytes := make([]byte, r.codec.ChatterIDBytes)
	_, err = r.reader.Read(chatterIDBytes)
	if err != nil && err != io.EOF {
		return nil, err
	}
	recChatterID := readUint(chatterIDBytes, r.codec.ChatterIDBytes)

	// Read in the bytes that make up the ChatterName field
	chatterNameLenBytes := make([]byte, r.codec.ChatterNameBytes)
	_, err = r.reader.Read(chatterNameLenBytes)
	if err != nil && err != io.EOF {
		return nil, err
	}
	chatterNameLen := readUint(chatterNameLenBytes, r.codec.ChatterNameBytes)
	if err == io.EOF && chatterNameLen != 0 {
		return nil, err
	}

	chatterNameByteArr := make([]byte, chatterNameLen)
	_, err = r.reader.Read(chatterNameByteArr)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// Read in the bytes that make up the Payload field
	payloadLenBytes := make([]byte, r.codec.PayloadBytes)
	_, err = r.reader.Read(payloadLenBytes)
	if err != nil && err != io.EOF {
		return nil, err
	}
	payloadLen := readUint(payloadLenBytes, r.codec.PayloadBytes)

	if err == io.EOF && payloadLen != 0 {
		return nil, err
	}

	payloadByteArr := make([]byte, payloadLen)
	_, err = r.reader.Read(payloadByteArr)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &Record {
		recType,
		recVersion,
		recRoomID,
		recChatterID,
		chatterNameByteArr,
		payloadByteArr,
	}, nil
}


func readUint(byteSli []byte, fieldSize ByteSize) uint {
	reader := bytes.NewReader(byteSli)
	switch fieldSize {
	case OneByte:
		var tempUInt uint8
		err := binary.Read(reader, binary.BigEndian, &tempUInt)
		if err != nil{
			fmt.Print("ERROR ONEBYTE: ", err)
		}
		return uint(tempUInt)
	case TwoBytes:
		var tempUInt uint16
		err := binary.Read(reader, binary.BigEndian, &tempUInt)
		if err != nil{
			fmt.Print("ERROR TWOBYTES: ", err)
		}
		return uint(tempUInt)
	case FourBytes:
		var tempUInt uint32
		err := binary.Read(reader, binary.BigEndian, &tempUInt)
		if err != nil{
			fmt.Print("ERROR FOURBYTES: ", err)
		}
		return uint(tempUInt)
	case EightBytes:
		var tempUInt uint64
		err := binary.Read(reader, binary.BigEndian, &tempUInt)
		if err != nil{
			fmt.Print("ERROR EIGHTBYTES: ", err)
		}
		return uint(tempUInt)
	default:
		return 0
	}
}


// Exported functions
func DecodeMsg(msg []byte) (Record, error) {

	var record Record
	reader := bytes.NewReader(msg)
	tlvReader := newReader(reader, RecordCodec)
	next, readingError := tlvReader.next()

	if readingError != nil {
		record = Record{
			Type: 0x8C,
			Version: 1,
			RoomID: 0,
			ChatterID: 0,
			ChatterName: []byte("ERROR"),
			Payload: []byte("ERROR"),}
	} else {
		record = *next
	}
	return record, readingError
}


func EncodeMsg(msgType uint, msg string, usrName string, roomNum uint) (*bytes.Buffer, error) {

	buf := new(bytes.Buffer)
	tlvWriter := newWriter(buf, RecordCodec)
	record := &Record{
			Type:        msgType,
			Version:     1,
			RoomID:      roomNum,
			ChatterID:   255,
			ChatterName: []byte(usrName),
			Payload:     []byte(msg),
		}

	err := tlvWriter.write(record)

	return buf, err
}