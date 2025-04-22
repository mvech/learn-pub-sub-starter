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

	gameState := gamelogic.NewGameState(username)

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			err = gameState.CommandSpawn(input)
			if err != nil {
				fmt.Printf("%s\n", err)
			}
		case "quit":
			gamelogic.PrintQuit()
			return
		case "move":
			_, err = gameState.CommandMove(input)
			if err != nil {
				fmt.Printf("%s\n", err)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Printf("Spamming not allowed yet!")
		default:
			fmt.Printf("Command not found! Try <help>\n")
		}
	}

}
