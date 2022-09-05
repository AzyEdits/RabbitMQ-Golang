package main

import (
	"bytes"
	"fmt"
	"helloworld/broker"
	"log"
	"time"

	"github.com/pkg/errors"
)

func main() {
	conn, ch, err := broker.RabbitMQ()

	if err != nil {
		panic(err)
	}

	defer func() {
		ch.Close()
		conn.Close()
	}()

	q, err := ch.QueueDeclare("task_queue", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err, "failed to declare queue")
	}

	err = ch.Qos(1, 0, false)

	if err != nil {
		log.Fatal(err, "failed to set Qos")
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)

	if err != nil {
		panic(errors.Wrap(err, "failed to consume queue"))
	}

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
