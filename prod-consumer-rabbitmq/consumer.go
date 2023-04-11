package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// printN prints the first N numbers
func printN(n uint32) {
	for i := uint32(0); i <= n; i++ {
		fmt.Printf("\r%d / %d", i, n)
	}
}

const QNAME = "Q"

func run() error {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("could not connect to rabbitmq: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("could not open channel: %w", err)
	}
	defer ch.Close()

	// We create a Queue to read messages from.
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

	// Tells RabbitMQ to not dispatch more than one message to a worker at a
	// time.
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return fmt.Errorf("could not set up fair dispatch: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("could not consume messages: %w", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %v\n", d.Body)
			payload := binary.LittleEndian.Uint32(d.Body)
			fmt.Printf("Decoded payload to %d\n", payload)

			// Simulate some working behaviour.
			time.Sleep(1 * time.Second)
			printN(payload)

			// Acknowledge the message so RabbitMQ removes it from the queue.
			d.Ack(false)
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("An error occured: %s", err.Error())
		os.Exit(1)
	}
}
