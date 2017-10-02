package rmq

import (
	"fmt"
	"os"
)

type rabbitMQConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}

// get formatted amqp address
func (c *rabbitMQConfig) getAddress() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", c.User, c.Password, c.Host, c.Port)
}

// get ENV variables (see https://hub.docker.com/_/rabbitmq/ for more information)
// and construct rabbitMQConfig object
func newRabbitConfig() *rabbitMQConfig {
	user := os.Getenv(ENV_RMQ_USER)
	password := os.Getenv(ENV_RMQ_PASSWORD)
	host := os.Getenv(ENV_RMQ_HOST)
	port := os.Getenv(ENV_RMQ_PORT)
	return &rabbitMQConfig{user, password, host, port}
}
