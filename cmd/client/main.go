package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ncvyn/peril-game/cmd"
	"github.com/ncvyn/peril-game/internal/gamelogic"
	"github.com/ncvyn/peril-game/internal/pubsub"
	"github.com/ncvyn/peril-game/internal/routing"
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
		handlerMove(gs, ch),
	)

	pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		"war",
		routing.WarRecognitionsPrefix+".*",
		pubsub.DurableQueue,
		handlerWar(gs, ch),
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
			if len(words) < 2 {
				fmt.Println("Usage: spam <number_of_moves>")
				continue
			}
			interations, err := strconv.Atoi(words[1])
			if err != nil {
				fmt.Println("Invalid number of moves:", err)
				continue
			}
			for i := range interations {
				err := pubsub.PublishGob(
					ch,
					routing.ExchangePerilTopic,
					routing.GameLogSlug+"."+gs.GetUsername(),
					routing.GameLog{
						CurrentTime: time.Now(),
						Message:     gamelogic.GetMaliciousLog(),
						Username:    gs.GetUsername(),
					},
				)
				if err != nil {
					fmt.Printf("Error publishing spam %d: %v\n", i+1, err)
				}
			}
		case "quit":
		case "exit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
