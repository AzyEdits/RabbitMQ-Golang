package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672")

	if err != nil {
		log.Fatal(err, "failed to connect to RabbitMQ")
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err, "failed to get channel")
	}

	defer ch.Close()

	q, err := ch.QueueDeclare("", false, false, true, false, nil)

	if err != nil {
		log.Fatal(err, "failed to declare queue")
	}

	err = ch.ExchangeDeclare("logs", amqp091.ExchangeFanout, true, false, false, false, nil)

	if err != nil {
		panic(errors.Wrap(err, "failed to declare exchange"))
	}

	err = ch.QueueBind(q.Name, "", "logs", false, nil)

	if err != nil {
		log.Fatal(err, "failed to bind queue")
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("[*] Waiting for messages. To Exit press CTRL+C")
	<-forever
}
