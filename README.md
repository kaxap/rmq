# rmq
A simple wrapper for Go RabbitMQ client http://github.com/streadway/amqp
Minimizes boilerplate code for using with dockerized RabbitMQ server

# example

```Go
package main

import (
	"github.com/kaxap/rmq"
	"log"
)

type something struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

func main() {

	// connect to the RabbitMQ server
	// connection parameters should be present as environment variables
	// i.e. RABBITMQ_DEFAULT_USER, RABBITMQ_DEFAULT_PASS, RABBITMQ_HOST, RABBITMQ_PORT
	// for more information see https://hub.docker.com/_/rabbitmq/
	queue, err := rmq.NewQueue("test_queue", false, 1, false, true, true)
	if err != nil {
		// could not connect or create channel
		log.Fatal(err)
	}
	defer queue.Close()

	var a something

	// consume a message
	msg, err := queue.Get(&a)
	log.Printf("message = %s\n", a.Data)

	// acknowledge the message
	msg.Ack(false)

	// modify data
	a.ID ++
	a.Data += " to you too"

	// send it back
	queue.Send(&a)
}

```

Now publish `{"id": 1, "data": "hello"}` to the "test_queue" and see what happens.
