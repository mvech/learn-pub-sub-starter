package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	state := routing.PlayingState{}
	state.IsPaused = true
	err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, state)
	if err != nil {
		log.Fatalf("unable to push message, %s", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Program is shutting down. Closing connections.")

}
