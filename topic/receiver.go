package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ.")
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open a channel.")
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
		return
	}

	err = ch.QueueBind(
		q.Name,
		"*.info", //another example's binding id : abc.# or #.id or *.xyz.* or *.xyz.id
		"logs_direct",
		false,
		nil,
	)
	if err != nil {
		failOnError(err, "Failed to Bind a Queue")
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		failOnError(err, "Failed to register a consumer")
		return
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			fmt.Println("Receieved a message:", string(d.Body))
		}
	}()

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
