package main

import (
	"fmt"

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

	pubsub.SubscribeGob(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.DurableQueue,
		handlerLogs(),
	)

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			fmt.Println("Pausing the game...")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
		case "resume":
			fmt.Println("Resuming the game...")
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
		case "quit":
		case "exit":
			fmt.Println("Quitting the server...")
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}
