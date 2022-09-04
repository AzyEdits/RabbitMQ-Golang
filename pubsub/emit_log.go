package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")

	if err != nil {
		panic(errors.Wrap(err, "failed to connect to RabbitMQ"))
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		panic(errors.Wrap(err, "failed to get channel"))
	}
	defer ch.Close()

	err = ch.ExchangeDeclare("logs", amqp091.ExchangeFanout, true, false, false, false, nil)

	if err != nil {
		panic(errors.Wrap(err, "failed to declare exchange"))
	}

	err = ch.Publish("logs", "", false, false, amqp091.Publishing{
		Headers:     map[string]interface{}{},
		ContentType: "text/plain",
		Body:        []byte(os.Args[1]),
	})

	if err != nil {
		panic(errors.Wrap(err, "failed to publish message"))
	}

	fmt.Println("Send message:", os.Args[1])
}
