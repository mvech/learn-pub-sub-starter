package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	url := "amqp://guest:guest@localhost:5672/"
	c, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("connection to amqp failed, %s", err)
	}
	defer c.Close()
	fmt.Println("Connection successful")

	username, err := gamelogic.ClientWelcome()

	_, _, err = pubsub.DeclareAndBind(c, routing.ExchangePerilDirect,
		routing.PauseKey+"."+username, routing.PauseKey, pubsub.Transient)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Program is shutting down. Closing connections.")
}
