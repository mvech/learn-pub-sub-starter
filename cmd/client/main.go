package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {

	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(mv gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(mv)
	}
}

func main() {
	fmt.Println("Starting Peril client...")
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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("error in client welcome, %s", err)
	}

	gameState := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		c,
		routing.ExchangePerilDirect,
		"pause."+gameState.GetUsername(),
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatalf("cannot subscribe to queue, %s", err)
	}

	err = pubsub.SubscribeJSON(
		c,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+gameState.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.Transient,
		handlerMove(gameState),
	)
	if err != nil {
		log.Fatalf("cannot subscribe to queue, %s", err)
	}

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
			move, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Printf("%s\n", err)
			}
			err = pubsub.PublishJSON(
				channel,
				routing.ExchangePerilTopic,
				routing.ArmyMovesPrefix+gameState.GetUsername(),
				move,
			)
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
