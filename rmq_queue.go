package rmq

import (
	"bytes"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

// Queue object with settings and message chan
type Queue struct {
	Channel       *Channel
	AmqpQueue     *amqp.Queue
	QueueName     string
	Durable       bool
	AutoDelete    bool
	Exclusive     bool
	NoWait        bool
	Arguments     amqp.Table
	PrefetchCount int
	PrefetchSize  int
	Global        bool
	Messages      <-chan amqp.Delivery
	AutoReconnect bool
	autoAck       bool
	consume       bool
}

// Queue object constructor
func NewQueue(name string, durable bool, prefetchCount int, autoAck, consume bool, autoReconnect bool) (*Queue, error) {

	queue := &Queue{
		nil,
		nil,
		name,
		durable,
		false,
		false,
		false,
		nil,
		prefetchCount,
		0,
		false,
		nil,
		autoReconnect,
		autoAck,
		consume,
	}

	return queue, queue.Connect()

}

// Queue object constructor, including queue arguments, such as map[string]interface{}{"x-queue-mode": "lazy"}
func NewQueueWithArgs(name string, durable bool, prefetchCount int, autoAck, consume bool, autoReconnect bool,
	args amqp.Table) (*Queue, error) {

	queue := &Queue{
		nil,
		nil,
		name,
		durable,
		false,
		false,
		false,
		args,
		prefetchCount,
		0,
		false,
		nil,
		autoReconnect,
		autoAck,
		consume,
	}

	return queue, queue.Connect()

}

// short syntax for NewQueue(name, durable=true, prefetchCount=0, autoAck=false, consume=false, autoReconnect=true)
func NewProducerQueue(name string) (*Queue, error) {
	return NewQueue(name, true, 0, false, false, true)
}

// short syntax for NewQueue(name, durable=true, prefetchCount=prefetchCount, autoAck=false, consume=true, autoReconnect=true)
func NewConsumerQueue(name string, prefetchCount int) (*Queue, error) {
	return NewQueue(name, true, prefetchCount, false, true, true)
}

// short syntax for producer queue with x-queue-mode: lazy args
func NewLazyProducerQueue(name string) (*Queue, error) {
	return NewQueueWithArgs(name, true, 0, false, false, true,
		map[string]interface{}{"x-queue-mode": "lazy"})
}

// short syntax for consumer queue with x-queue-mode: lazy args
func NewLazyConsumerQueue(name string, prefetchCount int) (*Queue, error) {
	return NewQueueWithArgs(name, true, prefetchCount, false, true, true,
		map[string]interface{}{"x-queue-mode": "lazy"})
}

func (q *Queue) Connect() error {

	// create a channel (and connection)
	ch, err := NewChannel()

	if err != nil {
		return err
	}

	q.Channel = ch

	// declare queue
	qrez, err := q.Channel.AmqpChannel.QueueDeclare(q.QueueName,
		q.Durable,
		q.AutoDelete,
		q.Exclusive,
		q.NoWait,
		q.Arguments)

	q.AmqpQueue = &qrez

	if err != nil {
		ch.Close()
		return err
	}

	if q.consume {
		// setup consumer settings

		err = ch.AmqpChannel.Qos(q.PrefetchCount, q.PrefetchSize, q.Global)
		if err != nil {
			ch.Close()
			return err
		}

		// bind queue to a chan
		messages, err := ch.AmqpChannel.Consume(q.QueueName,
			"",
			q.autoAck,
			false,
			false,
			false,
			nil)

		if err != nil {
			ch.Close()
			return err
		}

		q.Messages = messages
	}

	return nil
}

func (q *Queue) Reconnect() error {
	if q.Channel != nil {
		q.Channel.Close()
	}

	return q.Connect()
}

// Publish an amqp.Publishing structure
func (q *Queue) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {

	for {
		err1 := q.Channel.AmqpChannel.Publish(
			exchange,
			key,
			mandatory,
			immediate,
			msg,
		)

		if err1 != nil {
			if amqpErr, ok := err1.(*amqp.Error); ok {
				if amqpErr.Code == amqp.ChannelError {
					log.Println("Lost connection to queue manager.")
					if q.AutoReconnect {
						q.Reconnect()
						continue
					}
				}
			}
			return err1
		}

		break
	}

	return nil
}

// Publish an arbitrary structure by JSON serializing
func (q *Queue) PublishJSON(exchange, key string, mandatory, immediate bool, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return q.Publish(
		exchange,
		key,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
}

// Publish an arbitrary structure by JSON serializing to the default queue
func (q *Queue) Send(data interface{}) error {
	return q.PublishJSON("", q.QueueName, false, false, data)
}

// Consume a json message and deserialize it as JSON
// Note that if no AutoAck has been set for the queue, you must send ack manually:
// msg, err := queue.Get(&something)
// msg.Ack(false)
func (q *Queue) Get(v interface{}) (*amqp.Delivery, error) {
	msg := <-q.Messages

	if v != nil {
		d := json.NewDecoder(bytes.NewReader(msg.Body))
		d.UseNumber()
		err := d.Decode(v)

		if err != nil {
			log.Printf("Could not decode JSON: %s\n", err.Error())
		}
		return &msg, err
	}
	return &msg, nil
}

// Close associated channel
func (q *Queue) Close() error {
	return q.Channel.Close()
}
