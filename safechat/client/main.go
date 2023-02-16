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

type ConnState struct {
	pubKey *crypt.PublicKey
	symKey *[32]byte
}

func newState() ConnState {
	return ConnState{
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

func writeMsg(typ byte, msg string, s *ConnState) []byte {
	sends := []byte{typ}
	if s.symKey != nil && msg != "" {
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

	for {
		typ, msg := readMessage()
		sends := writeMsg(typ, msg, &state)
		_, err := connection.Write(sends)
		if err != nil {
			panic(err)
		}
		processMessage(connection, &state)
	}
}

func processMessage(connection net.Conn, s *ConnState) error {

	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		return err
	}
	if mLen == 0 {
		return errors.New("Received null message")
	}
	header := buffer[0]
	content := buffer[1:mLen]

	switch header {
	case SERVER_HELLO:
		fmt.Println("[server hello] received server hello")

		pubKey := &crypt.PublicKey{}
		pubKey.Unmarshal(content)

		s.pubKey = pubKey

		fmt.Printf("[server hello] public key is %+v\n", pubKey)

		symKey := generateSymKey()
		fmt.Printf("[server hello] generated sym key: %v\n", symKey)

		msg := pubKey.EncryptString(symKey[:])

		connection.Write(writeMsg(CLIENT_DONE, msg, s))

		s.symKey = &symKey

	case SERVER_MSG:
		fmt.Printf("[message] server encrypted message as: %s\n", buffer[1:mLen])

	case SERVER_DONE:
		fmt.Println("[server done] handshake complete")

	case ERROR:
		fmt.Printf("[error] received error: %s\n", content)

	default:
		fmt.Println("[error] handshake complete")
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
