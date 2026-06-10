package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/cmd"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	connection, err := cmd.Connect()
	if err != nil {
		panic(err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		panic(err)
	}

	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueType(pubsub.TransientQueue),
	)
	if err != nil {
		panic(err)
	}

	cmd.WaitForInterrupt()
}
