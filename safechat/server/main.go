// socket-server project main.go
package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

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

type AllConnections struct {
	m     *sync.Mutex
	conns map[net.Conn]ConnState
}

func NewAllConnections() AllConnections {
	return AllConnections{
		m:     &sync.Mutex{},
		conns: make(map[net.Conn]ConnState),
	}
}

func (a AllConnections) add(c net.Conn) {
	a.m.Lock()
	defer a.m.Unlock()
	a.conns[c] = NewConnState()
}

func (a AllConnections) del(c net.Conn) {
	a.m.Lock()
	defer a.m.Unlock()
	delete(a.conns, c)
}

func (a AllConnections) setPrivKey(c net.Conn, p crypt.PrivateKey) error {
	a.m.Lock()
	defer a.m.Unlock()
	if state, ok := a.conns[c]; ok {
		if state.priv != nil {
			return errors.New("private key was already set")
		}
		state.priv = &p
		a.conns[c] = state
	}
	return nil
}

func (a AllConnections) getPrivKey(c net.Conn) (crypt.PrivateKey, error) {
	a.m.Lock()
	defer a.m.Unlock()
	if _, ok := a.conns[c]; ok {
		return *a.conns[c].priv, nil
	}
	return crypt.PrivateKey{}, errors.New("connection not found")
}

func (a AllConnections) setSymKey(c net.Conn, sym [32]byte) error {
	a.m.Lock()
	defer a.m.Unlock()
	if state, ok := a.conns[c]; ok {
		if state.sym != nil {
			return errors.New("symmetric key was already set")
		}
		state.sym = &sym
		a.conns[c] = state
	}
	return nil
}

func (a AllConnections) getSymKey(c net.Conn) ([32]byte, error) {
	a.m.Lock()
	defer a.m.Unlock()
	if _, ok := a.conns[c]; ok {
		return *a.conns[c].sym, nil
	}
	return [32]byte{}, errors.New("connection not found")
}

func run() error {
	fmt.Println("Server Running...")

	allCons := NewAllConnections()

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
		if err != nil {
			fmt.Println("Error accepting client: ", err.Error())
		}
		fmt.Println("client connected")
		go processClient(connection, allCons)
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

func processClient(connection net.Conn, allCons AllConnections) {

	allCons.add(connection)
	defer allCons.del(connection)
	defer func() {
		fmt.Println("client disconnected")
	}()

	for {
		err := processMessage(connection, allCons)
		if err != nil {
			break
		}
	}
}

func processMessage(connection net.Conn, allCons AllConnections) error {
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
	case CLIENT_HELLO:
		pub, priv := crypt.GenerateKeyPair()
		err := allCons.setPrivKey(connection, priv)
		if err != nil {
			connection.Write(getMsg(ERROR, err.Error()))
			fmt.Println("[server log] received hello request twice")
			break
		}
		sends := []byte{SERVER_HELLO}

		pubBytes := pub.Marshal()
		sends = append(sends, pubBytes...)

		connection.Write(sends)

	case CLIENT_DONE:
		symKeyEncrypted := buffer[1:mLen]
		fmt.Printf("Received encrypted symKey: %v\n", symKeyEncrypted)

		privKey, _ := allCons.getPrivKey(connection)
		symKey := privKey.DecryptString(fmt.Sprintf("%s", symKeyEncrypted))
		fmt.Printf("Decrypted symKey: %v\n", symKey)

		symKey32 := [32]byte{}
		copy(symKey32[:], symKey[:])

		allCons.setSymKey(connection, symKey32)

	case CLIENT_MSG:
		fmt.Printf("[info] received encrypted message: %s\n", buffer[1:mLen])
		symkey, _ := allCons.getSymKey(connection)
		msg := crypt.DecryptAES(symkey[:], buffer[1:mLen])
		fmt.Printf("[info] decrypted message: %s\n", msg)
	default:
		return errors.New("Received null message")
	}
	return nil
}

func getMsg(typ byte, msg string) []byte {
	sends := []byte{typ}
	sends = append(sends, []byte(msg)...)
	return sends
}
