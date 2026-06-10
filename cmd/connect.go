package cmd

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp.Connection, error) {
	fmt.Println("Connecting to RabbitMQ server...")

	const url = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ server @", url)
	return conn, nil
}
