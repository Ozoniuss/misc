Producer-Consumer queue
-----------------------

This is a basic producer-consumer queue implemented with RabbitMQ. It provides the following features:

- Running a producer. This will keep asking for user input until canceled. The input is expected to be an unsigned integer that fits on 32 bits.
- Running a worker. This will keep reading from the queue. It will start printing on the same line all numbers between 0 and the number read from the queue (e.g. `123 / 5000`) until all numbers are completed.
- Fair dispatching. The RabbitMQ broker will dispatch messages to any available worker, of the first available worker. It does not use the default round-robin approach.
- Worker failsafe. This makes use of RabbitMQ's message acknowlledgment mechanism to ensure that if a worker crashes while processing a message, the message will still remain in the queue until it's processed by a different worker.
- Queue durability. If the RabbitMQ service crashes or quits, it will make the best effort to keep the messages that are in the queue (stronger persistence may be implemented but yeah...) 

Run RabbitMQ with:

```sh
docker compose up 
```

Run the producer with:

```sh
go run producer.go
```

Run any number of consumers with:

```sh
go run consumer.go
```