package main

import (
	"github.com/streadway/amqp"
)

// NewDeviceRPC calls AMQP to register new device
func NewDeviceRPC() (res string) {
	// Dialing RabbitMQ broker and creating a connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connecto to RabbitMQ")
	defer conn.Close()

	// Creating a new connection
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declaring consume queue for RPC callback
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	// Consuming on that queue waiting for response
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer for callbacks")

	// Creating a correlation ID to identify response
	corrID := randomString(32)

	// Call the RPC server registering a new device
	err = ch.Publish(
		"",
		"new_device_rpc",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       q.Name,
			Body:          []byte("Hola!"),
		})
	failOnError(err, "Failed to publish RPC message")

	// Wait for the response (correlation ID matched) and process it
	for d := range msgs {
		if corrID == d.CorrelationId {
			res = string(d.Body)
			break
		}
	}

	return
}

// NewReading calls the HistoricRegistry Service to save a new reading
func NewReading() (res string) {
	return "Pass"
}
