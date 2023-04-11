package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const QNAME = "Q"

func run() error {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		return fmt.Errorf("could not connect to rabbitmq: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("could not open channel: %w", err)
	}
	defer ch.Close()

	// We create a Queue to send the message to.
	q, err := ch.QueueDeclare(
		QNAME, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("could not declare queue: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; ; i++ {
		var num uint32

		// Be careful to read a number less than 2^32 so it doesn't fuck up
		// scan.
		fmt.Print("get number: ")
		_, err := fmt.Scan(&num)
		if err != nil {
			return fmt.Errorf("could not read number: %w", err)
		}

		if num == 0 {
			break
		}

		var payload = make([]byte, 4, 4)
		binary.LittleEndian.PutUint32(payload, num)

		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         payload,
			})
		// If there is an error publishing the message, a log will be displayed in the terminal.
		if err != nil {
			return fmt.Errorf("could not publish message: %w", err)
		}
		fmt.Printf("[x] Sent payload: %v\n", payload)
	}
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("An error occured: %s", err.Error())
		os.Exit(1)
	}
}
