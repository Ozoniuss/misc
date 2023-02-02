package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"

	crypt "safechat/encryption"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9987"
	SERVER_TYPE = "tcp"
)

const (
	CLIENT_HELLO byte = 0
	SERVER_HELLO byte = 1
	CLIENT_DONE  byte = 2
	SERVER_DONE  byte = 3
	ERROR        byte = 4
	CLIENT_MSG   byte = 5
	SERVER_MSG   byte = 6
	CLIENT_CLOSE byte = 7
	SERVER_CLOSE byte = 8
)

type state struct {
	pubKey *crypt.PublicKey
	symKey *[32]byte
}

func newState() state {
	return state{
		pubKey: nil,
		symKey: nil,
	}
}

func readMessage() (byte, string) {
	var typ byte
	var msg string
	fmt.Scanf("%d:%s\n", &typ, &msg)
	return typ, msg
}

func getMsg(typ byte, msg string, s state) []byte {
	sends := []byte{typ}
	if s.symKey != nil {
		fmt.Println("sddas")
		msg = fmt.Sprintf("%s", crypt.EncryptAES(s.symKey[:], []byte(msg)))
	}
	sends = append(sends, []byte(msg)...)
	return sends
}

func main() {
	//establish connection
	connection, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
	}

	state := newState()

	//send some data
	for {
		typ, msg := readMessage()
		sends := getMsg(typ, msg, state)
		_, err := connection.Write(sends)
		if err != nil {
			panic(err)
		}
		processMessage(connection, &state)
	}

	// mLen, err := connection.Read(buffer)
	// if err != nil {
	// 	fmt.Println("Error reading:", err.Error())
	// }
	// fmt.Println("Received: ", string(buffer[:mLen]))
	// defer connection.Close()
}

func processMessage(connection net.Conn, s *state) error {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		return err
	}
	if mLen == 0 {
		return errors.New("Received null message")
	}
	header := buffer[0]
	switch header {
	case ERROR:
		fmt.Printf("[error] %s\n", string(buffer[1:mLen]))
	case SERVER_HELLO:
		fmt.Println("[info] received server hello")

		pubKey := &crypt.PublicKey{}
		pubKey.Unmarshal(buffer[1:mLen])

		s.pubKey = pubKey

		fmt.Printf("public key is %+v\n", pubKey)

		symKey := generateSymKey()
		fmt.Printf("generated sym key: %v\n", symKey)

		msg := pubKey.EncryptString(symKey[:])

		connection.Write(getMsg(CLIENT_DONE, msg, *s))

		s.symKey = &symKey

	case SERVER_MSG:
		fmt.Printf("[info] received encrypted message: %s\n", buffer[1:mLen])
		//fmt.Printf("[info] decrypted message: %s\n")

	default:
		return errors.New("Received message with invalid header")
	}
	return nil
}

func generateSymKey() [32]byte {
	var key32 [32]byte
	key := make([]byte, 32)
	rand.Read(key)
	copy(key32[:], key[:])
	return key32
}
