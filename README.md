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

	// this is short syntax for a durable consumer queue,
	// if you need to create a non-durable queue, please use NewQueue constructor (see "Constructors" chapter below)
	inputQueue, err := rmq.NewConsumerQueue("input_queue", 1)
	if err != nil {
		// could not connect or create channel
		log.Fatal(err)
	}
	defer inputQueue.Close()

	// this is short syntax for a durable producer queue
	outputQueue, err := rmq.NewProducerQueue("output_queue", 1)
	if err != nil {
		// could not connect or create channel
		log.Fatal(err)
	}
	defer outputQueue.Close()

	var a something

	// consume a json encoded message
	msg, err := inputQueue.Get(&a)
	log.Printf("message = %s\n", a.Data)

	// acknowledge the message
	msg.Ack(false)

	// modify data
	a.ID ++
	a.Data += " to you too"

	// send it to output queue
	outputQueue.Send(&a)
}

```

Now publish `{"id": 1, "data": "hello"}` to the "test_queue" and see what happens.


# Even shorter syntax if you don't need to send messages very often

```golang
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

    const maxRetry := 10
    rmq.SendAndForget("some_queue_where_I_want_to_dump_something", &something{ID: 1, Data: "hello!"}, maxRetry)
    // this will try to connect to a durable queue with the given name and send the data in "something" struct
    // if connection or sending failed it will retry up to 10 times. Then it will close the connection

```

Note if you need to send messages to non durable queues, please use `rmq.SendAndForgetNonDurable`. There is also
`SendAndForgetLazy` available for lazy queues.


# Constructors

There is 6 types of queue constructors:

```golang
func NewQueue(
    name string, // queue name
    durable bool, // durable flag
    prefetchCount int, // how many message to prefetch for this client
    autoAck bool, // auto_ack flag
    consume bool, // true for consumer/producer worker, false for producer-only worker
    autoReconnect bool, // true for auto reconnect to rabbitmq server
) (*Queue, error)
```

```golang
func NewQueueWithArgs(
    name string, // queue name
    durable bool, // durable flag
    prefetchCount int, // how many message to prefetch for this client
    autoAck bool, // auto_ack flag
    consume bool, // true for consumer/producer worker, false for producer-only worker
    autoReconnect bool, // true for auto reconnect to rabbitmq server
    args amqp.Table, // map[string]interface{} of queue arguments
) (*Queue, error)
```

```golang
NewProducerQueue(name string) (*Queue, error)
// short syntax for NewQueue(name, durable=true, prefetchCount=0, autoAck=false, consume=false, autoReconnect=true)
```

```golang
NewConsumerQueue(name string, prefetchCount int) (*Queue, error)
// short syntax for NewQueue(name, durable=true, prefetchCount=prefetchCount, autoAck=false, consume=true, autoReconnect=true)
```

```golang
NewLazyProducerQueue(name string) (*Queue, error)
// short syntax for producer queue with x-queue-mode: lazy args
```

```golang
NewLazyConsumerQueue(name string, prefetchCount int) (*Queue, error)
// short syntax for consumer queue with x-queue-mode: lazy args
```

Note that lazy queues are often used when a queue is expected to be frequently flooded. In lazy mode RabbitMQ pages out the messages on disk when possible.
For more information see https://www.rabbitmq.com/lazy-queues.html
