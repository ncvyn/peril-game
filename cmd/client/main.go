package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/cmd"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	conn, err := cmd.Connect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		panic(err)
	}

	gs := gamelogic.NewGameState(username)

	pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.TransientQueue,
		handlerPause(gs),
	)

	pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		"army_moves."+username,
		"army_moves.*",
		pubsub.TransientQueue,
		handlerMove(gs),
	)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err := gs.CommandSpawn(words)
			if err != nil {
				fmt.Println("Error spawning unit:", err)
			}
		case "move":
			mv, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println("Error moving unit:", err)
			}
			pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				"army_moves."+gs.GetUsername(),
				mv,
			)
			fmt.Println("Move command sent!")

		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
		case "exit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
