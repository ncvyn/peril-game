package cmd

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Connection, error) {
	fmt.Println("Connecting to RabbitMQ server...")

	const url = "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})

	fmt.Println("Successfully connected to RabbitMQ server @", url)
	return connection, nil
}
