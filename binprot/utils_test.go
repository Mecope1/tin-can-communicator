package binprot

import (
	"fmt"
	"testing"
)


func TestEncodeMsg(t *testing.T) {
	ans, err := EncodeMsg(0x8A, "Hello, World!", "user1", 2)
	ansByteSli := ans.Bytes()

	fmt.Println("ansByteSli",ansByteSli)
	encodedByteSli := []byte{
		138, 1, 2, 255, 0, 5, 117, 115, 101, 114, 49, 0, 0, 0, 13, 72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33,
		}
	if err != nil {
		t.Errorf("EncodeMsg(0x8A, \"Hello, World!\", \"user1\", 2) failed encoding")
	}

	if len(ansByteSli) != len(encodedByteSli) {
		t.Errorf("EncodeMsg(0x8A, \"Hello, World!\", \"user1\", 2) got length of %d expected %d", len(ansByteSli), len(encodedByteSli))
	} else {

		for ind, val := range ansByteSli {
			if val != encodedByteSli[ind] {
				t.Errorf("EncodeMsg(0x8A, \"Hello, World!\", \"user1\", 2) byte#%d = %s expected %s", ind, string(val), string(ansByteSli[ind]))
			}
		}

	}

}

func TestDecodeMsg(t *testing.T) {

	encodedByteSli := []byte{
		138, 1, 2, 255, 0, 5, 117, 115, 101, 114, 49, 0, 0, 0, 13, 72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33,
	}

	// Contains 402 characters
	//longMessage402 := "qwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertqwertKAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOrange"
	// ans, err := EncodeMsg(0x8A, longMessage402, "user1", 2)

	ans, err := EncodeMsg(0x8A, "Hello, World!", "user1", 2)

	fmt.Println("EncodedMSG Result", ans.Bytes())
	fmt.Println("Hand-encoded Result", encodedByteSli)
	//msg, err := DecodeMsg(ans.Bytes())
	msg, err := DecodeMsg(encodedByteSli)

	if err != nil {
		t.Errorf("DecodeMsg(msg) failed encoding")
	} else if fmt.Sprintf("%T", msg) != "binprot.Record" {
		t.Errorf("EncodeMsg(msg) produced %s expected binprot.Record", fmt.Sprintf("%T", msg) )
	} else {
		// Checks fields from the decoded message
		switch {
		case msg.Type != 0x8A:
			t.Errorf("EncodeMsg(msg) produced %d expected 0x8A", msg.Type)
		case msg.Version != 1:
			t.Errorf("EncodeMsg(msg) produced %d expected 1", msg.Version)
		case msg.RoomID != 2:
			t.Errorf("EncodeMsg(msg) produced %d expected 2", msg.RoomID)
		case msg.ChatterID != 255:
			t.Errorf("EncodeMsg(msg) produced %d expected 255", msg.ChatterID)
		case string(msg.ChatterName) != "user1":
			t.Errorf("EncodeMsg(msg) produced %s expected user1", msg.ChatterName)
		//case string(msg.Payload) != longMessage402:
		//	t.Errorf("EncodeMsg(msg) produced %s expected %s", string(msg.Payload), longMessage402)
		case string(msg.Payload) != "Hello, World!":
			t.Errorf("EncodeMsg(msg) produced %s expected Hello, World!", string(msg.Payload))
		}
		println("msg.Payload", string(msg.Payload))

	}

}