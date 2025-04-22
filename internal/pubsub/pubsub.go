package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {

	valBytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("unable to mashal  into json, %s", err)
	}

	publishing := amqp.Publishing{}
	publishing.ContentType = "application/json"
	publishing.Body = valBytes
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, publishing)
	if err != nil {
		return fmt.Errorf("unable to publish, %s", err)
	}

	return nil
}

const (
	Durable int = iota
	Transient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error creating a channel, %s", err)
	}

	queue, err := channel.QueueDeclare(queueName,
		simpleQueueType == Durable,
		simpleQueueType == Transient,
		simpleQueueType == Transient, false, nil)

	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error creating queue, %s", err)
	}

	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("error binding queue, %s", err)
	}

	return channel, queue, nil

}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
	handler func(T),
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
			handler(*m)
			d.Ack(false)
		}
	}()
	return nil
}
