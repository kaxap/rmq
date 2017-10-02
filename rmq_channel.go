package rmq

import (
	"github.com/streadway/amqp"
	"log"
)

// Channel contains an active connection to RabbitMQ and a reference to AMQP Channel object
type Channel struct {
	Conn        *amqp.Connection
	AmqpChannel *amqp.Channel
}

// Channel constructor
// Creates connection from ENV parameters (see consts.go)
func NewChannel() (*Channel, error) {
	config := newRabbitConfig()
	addr := config.getAddress()

	log.Printf("RabbitMQ: Connecting to %s:%s...\n", config.Host, config.Port)
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Channel{Conn:conn, AmqpChannel:ch}, nil
}

// Close channel and connection to RabbitMQ
func (c *Channel) Close() error {
	err := c.AmqpChannel.Close()
	if err != nil {
		return err
	}

	return c.Conn.Close()
}