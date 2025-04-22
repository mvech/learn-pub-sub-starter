package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	url := "amqp://guest:guest@localhost:5672/"
	c, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("connection to amqp failed, %s", err)
	}
	defer c.Close()

	fmt.Println("Connection successful")

	channel, err := c.Channel()
	if err != nil {
		log.Fatalf("error in creating channel, %s", err)
	}

	_, _, err = pubsub.DeclareAndBind(c, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.Durable)
	if err != nil {
		log.Fatalf("Error creating and binding queue, %s", err)
	}

	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			state := routing.PlayingState{}
			state.IsPaused = true
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, state)
			if err != nil {
				log.Fatalf("unable to push message, %s", err)
			}
			fmt.Print("Sending pause message")
		case "resume":
			state := routing.PlayingState{}
			state.IsPaused = false
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, state)
			if err != nil {
				log.Fatalf("unable to push message, %s", err)
			}
			fmt.Println("Sending resume message")
		case "quit":
			fmt.Println("Exiting game")
			return
		}
	}

}
