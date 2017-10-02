package rmq

import (
	"github.com/streadway/amqp"
	"encoding/json"
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
}

// Queue object constructor
func NewQueue(name string, durable bool, prefetchCount int, autoAck, consume bool) (*Queue, error) {

	// create a channel (and connection)
	ch, err := NewChannel()

	if err != nil {
		return nil, err
	}

	queue := &Queue{
		ch,
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
	}

	// declare queue
	q, err := ch.AmqpChannel.QueueDeclare(name,
		queue.Durable,
		queue.AutoDelete,
		queue.Exclusive,
		queue.NoWait,
		queue.Arguments)

	queue.AmqpQueue = &q

	if err != nil {
		ch.Close()
		return nil, err
	}

	if consume {
		// setup consumer settings

		err = ch.AmqpChannel.Qos(queue.PrefetchCount, queue.PrefetchSize, queue.Global)
		if err != nil {
			ch.Close()
			return nil, err
		}

		// bind queue to a chan
		messages, err := ch.AmqpChannel.Consume(name,
			"",
			autoAck,
			false,
			false,
			false,
			nil)

		if err != nil {
			ch.Close()
			return nil, err
		}

		queue.Messages = messages
	}

	return queue, nil

}

// Publish an amqp.Publishing structure
func (q *Queue) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return q.Channel.AmqpChannel.Publish(
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
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
	msg := <- q.Messages
	err := json.Unmarshal(msg.Body, v)
	return &msg, err
}

// Close associated channel
func (q *Queue) Close() error {
	return q.Channel.Close()
}