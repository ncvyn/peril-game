package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	unmarshaller := func(body []byte) (T, error) {
		var msg T
		err := json.Unmarshal(body, &msg)
		return msg, err
	}
	return subscribe(conn, exchange, queueName, key, queueType, handler, unmarshaller)
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
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
			var err error
			msg, err = unmarshaller(delivery.Body)
			if err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			ack := handler(msg)
			switch ack {
			case Ack:
				if err := delivery.Ack(false); err != nil {
					log.Printf("Error acknowledging message: %v", err)
				}
			case NackRequeue:
				if err := delivery.Nack(false, true); err != nil {
					log.Printf("Error negatively acknowledging and requeuing message: %v", err)
				}
			case NackDiscard:
				if err := delivery.Nack(false, false); err != nil {
					log.Printf("Error negatively acknowledging and discarding message: %v", err)
				}
			}
		}
	}()

	return nil
}
