package rmq

import (
	"github.com/pkg/errors"
	"log"
	"time"
)

// connects to rabbitMQ as a producer, sends the data and disconnects
func SendAndForget(queueName string, data interface{}, retries int) {
	if retries == 0 {
		return
	}
	q, err := NewProducerQueue(queueName)
	if err != nil {
		// wait 5 seconds and try again
		log.Println(errors.Wrapf(err, "RMQ SendAndForget: Could not send to queue %s, retries left %d", queueName, retries))
		<-time.After(5 * time.Second)
		SendAndForget(queueName, data, retries-1)
	}
	defer q.Close()

	err = q.Send(data)
	if err != nil {
		// wait 5 seconds and try again
		log.Println(errors.Wrapf(err, "RMQ SendAndForget: Could not send to queue %s, retries left %d", queueName, retries))
		<-time.After(5 * time.Second)
		SendAndForget(queueName, data, retries-1)
	}

}

// connects to rabbitMQ as a producer, sends the data and disconnects
func SendAndForgetLazy(queueName string, data interface{}, retries int) {
	if retries == 0 {
		return
	}
	q, err := NewLazyProducerQueue(queueName)
	if err != nil {
		// wait 5 seconds and try again
		log.Println(errors.Wrapf(err, "RMQ SendAndForget: Could not send to queue %s, retries left %d", queueName, retries))
		<-time.After(5 * time.Second)
		SendAndForgetLazy(queueName, data, retries-1)
	}
	defer q.Close()

	err = q.Send(data)
	if err != nil {
		// wait 5 seconds and try again
		log.Println(errors.Wrapf(err, "RMQ SendAndForget: Could not send to queue %s, retries left %d", queueName, retries))
		<-time.After(5 * time.Second)
		SendAndForgetLazy(queueName, data, retries-1)
	}

}

// connects to rabbitMQ as a producer, sends the data and disconnects
func SendAndForgetNonDurable(queueName string, data interface{}, retries int) {
	if retries == 0 {
		return
	}
	q, err := NewQueue(queueName, false, 0, false, false, true)
	if err != nil {
		// wait 5 seconds and try again
		log.Println(errors.Wrapf(err, "RMQ SendAndForget: Could not send to queue %s, retries left %d", queueName, retries))
		<-time.After(5 * time.Second)
		SendAndForgetNonDurable(queueName, data, retries-1)
	}
	defer q.Close()

	err = q.Send(data)
	if err != nil {
		// wait 5 seconds and try again
		log.Println(errors.Wrapf(err, "RMQ SendAndForget: Could not send to queue %s, retries left %d", queueName, retries))
		<-time.After(5 * time.Second)
		SendAndForgetNonDurable(queueName, data, retries-1)
	}

}
