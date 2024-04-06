package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"outbox/articles"
	"time"
)

func run() error {
	lis, err := net.Listen("tcp", "127.0.0.1:13311")
	if err != nil {
		return fmt.Errorf("could not listen: %s", err.Error())
	}
	fmt.Println("starting consumer...")

	conn, err := lis.Accept()
	if err != nil {
		return fmt.Errorf("could not accept connection: %s", err.Error())
	}
	fmt.Println("accepted connection")

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	for {
		var newArticles = make([]articles.ArticleLikedEvent, 0)
		err := dec.Decode(&newArticles)
		if err != nil {
			fmt.Printf("failed decoding article liked events: %s", err.Error())
			continue
		}
		fmt.Println(newArticles, len(newArticles))
		time.Sleep(15 * time.Second)

		// ack
		err = enc.Encode(true)
		if err != nil {
			fmt.Printf("failed sending ack: %s", err.Error())
		}
	}
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("error running consumer: %s\n", err.Error())
		os.Exit(1)
	}
}
