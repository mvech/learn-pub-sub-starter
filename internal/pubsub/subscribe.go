package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {

	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}

	deliveries, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		defer channel.Close()
		for d := range deliveries {
			m := new(T)
			err = json.Unmarshal(d.Body, m)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			ack := handler(*m)
			switch ack {
			case Ack:
				d.Ack(false)
				fmt.Println("Ack")
			case NackRequeue:
				d.Nack(false, true)
				fmt.Println("NackRequeue")
			case NackDiscard:
				d.Nack(false, false)
				fmt.Println("NackDiscard")
			}
		}
	}()
	return nil
}
