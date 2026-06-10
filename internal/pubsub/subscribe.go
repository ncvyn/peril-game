package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {
	ch, q, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		return err
	}
	deliveryCh, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range deliveryCh {
			var msg T
			if err := json.Unmarshal(delivery.Body, &msg); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			handler(msg)
			delivery.Ack(false)
		}
	}()

	return nil
}
