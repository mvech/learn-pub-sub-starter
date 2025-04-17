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
