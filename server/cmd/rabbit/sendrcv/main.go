package main

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	//channel 虚拟的连接
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	q, err := ch.QueueDeclare(
		"go_q1",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	go consume("c1",conn, q.Name)
	go consume("c2",conn, q.Name)
	i := 0
	for {
		i++
		err := ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				Body: []byte(fmt.Sprintf("message %d", i)),
			},
		)
		if err != nil {
			fmt.Println(err.Error())
		}
		time.Sleep(200 * time.Millisecond)
	}

}

func consume(consumer string,conn *amqp.Connection, q string) {
	ch, err := conn.Channel()

	if err != nil {
		panic(err)
	}
	defer ch.Close()
	msgs, err := ch.Consume(
		q,
		consumer,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		fmt.Printf("%s: %s\n",consumer, msg.Body)
	}

}
