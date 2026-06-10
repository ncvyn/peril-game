package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	const url = "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})

	fmt.Println("Successfully connected to RabbitMQ server @", url)
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
	waitForInterrupt()
}

func waitForInterrupt() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	fmt.Println( /*^C*/ "opy that, shutting down the client...")
}
