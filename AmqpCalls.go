package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// AmqpClient holds a connection to the broker
type AmqpClient struct {
	conn *amqp.Connection
}

// ConnectToBroker sets up a connection to an AMQP Broker
func (a *AmqpClient) ConnectToBroker(connectionString string) {
	if connectionString == "" {
		panic("Connection string not set. Can't connecto to the broker")
	}

	var err error
	a.conn, err = amqp.Dial(fmt.Sprintf("%s/", connectionString))
	if err != nil {
		failOnError(err, "Failed to connect to broker: "+connectionString)
	}
}

// NewDeviceRPC calls AMQP to register new device
func (a *AmqpClient) NewDeviceRPC(jsonObj newRegister) (res string, code int) {
	// Creating a new connection
	ch, err := a.conn.Channel()
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
	jsonBytes, err := json.Marshal(jsonObj)
	err = ch.Publish(
		"",
		"device_management_rpc",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       q.Name,
			Body:          jsonBytes,
		})
	failOnError(err, "Failed to publish RPC message")

	// Wait for the response (correlation ID matched) and process it
	for msg := range msgs {
		if corrID == msg.CorrelationId {
			headers := msg.Headers
			res = string(msg.Body)
			codeUint, ok := headers["status"].(uint8)
			log.Println(ok, codeUint)
			if !ok {
				code = 0
			} else {
				code = int(codeUint)
			}
			break
		}
	}

	return
}

// NewReading calls the HistoricRegistry Service to save a new reading
func (a *AmqpClient) NewReading() (res string) {
	return "Pass"
}

// Close closes the connection to the broker
func (a *AmqpClient) Close() {
	if a.conn != nil {
		log.Printf("%s", "Closing AMQP connection")
		a.conn.Close()
	}
}
