package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

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

	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err, "failed to declare queue")
	}

	err = ch.Qos(1, 0, false)

	if err != nil {
		log.Fatal(err, "failed to set Qos")
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Println("Done")
			d.Ack(false) //Permite que solo reciba un mensaje por el canal dejando los demás mensajes encolados
		}
	}()

	fmt.Println("[*] Waiting for messages. To Exit press CTRL+C")
	<-forever
}
