package main

import (
	"fmt"

	"github.com/ncvyn/peril-game/internal/gamelogic"
	"github.com/ncvyn/peril-game/internal/pubsub"
	"github.com/ncvyn/peril-game/internal/routing"
)

func handlerLogs() func(routing.GameLog) pubsub.AckType {
	return func(gl routing.GameLog) pubsub.AckType {
		defer fmt.Print("> ")
		err := gamelogic.WriteLog(gl)
		if err != nil {
			fmt.Printf("Error writing game log: %v\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
