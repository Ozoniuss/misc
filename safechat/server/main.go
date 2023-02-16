// socket-server project main.go
package main

import (
	"errors"
	"fmt"
	"net"
	"os"

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

// ConnState represents the state of the connection with the client.
type ConnState struct {
	clientHello bool
	priv        *crypt.PrivateKey
	sym         *[32]byte
}

func NewConnState() ConnState {
	return ConnState{
		clientHello: false,
		priv:        nil,
		sym:         nil,
	}
}

func (state *ConnState) setPrivKey(p crypt.PrivateKey) error {
	if state.priv != nil {
		return errors.New("private key was already set")
	}
	state.priv = &p
	return nil
}

func (state *ConnState) getPrivKey() crypt.PrivateKey {
	return *state.priv
}

func (state *ConnState) setSymKey(s [32]byte) error {
	if state.sym != nil {
		return errors.New("symmetric key was already set")
	}
	state.sym = &s
	return nil
}

func (state *ConnState) getSymKey() [32]byte {
	return *state.sym
}

func run() error {
	fmt.Println("Server Running...")

	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	defer server.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	fmt.Println("Waiting for client...")

	for {
		connection, err := server.Accept()
		state := NewConnState()
		if err != nil {
			fmt.Println("Error accepting client: ", err.Error())
		}
		fmt.Println("client connected")
		processClient(connection, &state)
	}
}

func main() {
	// Running the code in a separate function allows executing the deferred
	// functions before exiting with code 1. The call os.Exit() stops the
	// subsequent deferred functions.
	err := run()
	if err != nil {
		fmt.Printf("An error occured: %s", err.Error())
		os.Exit(1)
	}
}

func processClient(connection net.Conn, state *ConnState) {

	defer func() {
		fmt.Println("client disconnected")
	}()

	for {
		err := processMessage(connection, state)
		if err != nil {
			break
		}
	}
}

func processMessage(connection net.Conn, state *ConnState) error {
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
	case CLIENT_HELLO:
		fmt.Println("[client hello]: received client hello")
		pub, priv := crypt.GenerateKeyPair()
		err := state.setPrivKey(priv)
		if err != nil {
			connection.Write(writeMsg(ERROR, "client hello failed: received hello request twice"))
			fmt.Println("[server log] received hello request twice")
			break
		}

		pubBytes := pub.Marshal()
		sends := writeMsg(SERVER_HELLO, string(pubBytes))

		connection.Write(sends)

	case CLIENT_DONE:
		// At this step it is assumed that the client returned his symmetric
		// key.
		symKeyEncrypted := content
		fmt.Printf("[client done] received encrypted symmetric key: %v\n", symKeyEncrypted)

		privKey := state.getPrivKey()
		symKey := privKey.DecryptString(fmt.Sprintf("%s", symKeyEncrypted))
		fmt.Printf("[client done] decrypted symmetrick key is: %v\n", symKey)

		symKey32 := [32]byte{}
		copy(symKey32[:], symKey[:])

		state.setSymKey(symKey32)

		sends := writeMsg(SERVER_DONE, "")
		connection.Write(sends)

	case CLIENT_MSG:
		fmt.Printf("[message] received encrypted message: %s\n", content)
		symkey := state.getSymKey()
		msg := crypt.DecryptAES(symkey[:], buffer[1:mLen])
		fmt.Printf("[message] decrypted message: %s\n", msg)

		sends := []byte{SERVER_MSG}
		sends = append(sends, content...)

		connection.Write(sends)

	default:
		fmt.Printf("[error] received invalid header")
		sends := writeMsg(ERROR, "received invalid header")
		connection.Write(sends)
	}
	return nil
}

func writeMsg(typ byte, msg string) []byte {
	sends := []byte{typ}
	if msg != "" {
		sends = append(sends, []byte(msg)...)
	}
	return sends
}
