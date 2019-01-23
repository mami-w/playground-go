package service

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type dispatcher struct {
	channel       *amqp.Channel
	queueName     string
	mandatorySend bool
}

func newDispatcher(channel *amqp.Channel, queueName string, mandatorySend bool) (disp *dispatcher) {
	disp = &dispatcher{channel:channel, queueName:queueName, mandatorySend:mandatorySend}
	return disp
}

func (q *dispatcher) DispatchMessage(message interface{}) (err error) {
	fmt.Printf("Dispatching message to queue %s\n", q.queueName)
	body, err := json.Marshal(message)
	if err == nil {
		err = q.channel.Publish(
			"",              // exchange
			q.queueName,     // routing key
			q.mandatorySend, // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			fmt.Printf("Failed to dispatch message: %s\n", err)
		}
	} else {
		fmt.Printf("Failed to marshal message %v (%s)\n", message, err)
	}
	return
}