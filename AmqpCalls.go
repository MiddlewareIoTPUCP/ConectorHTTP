package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/houseofcat/turbocookedrabbit/v2/pkg/tcr"
	"github.com/streadway/amqp"
)

// ConnectToBroker sets up a connection to an AMQP Broker
func ConnectToBroker(connectionString string) *tcr.ConnectionPool {
	if connectionString == "" {
		panic("Connection string not set. Can't connecto to the broker")
	}

	// Create connection config from TCR library to get a connection pool
	config := &tcr.PoolConfig{
		ConnectionName:       "ConectorHTTP",
		URI:                  connectionString,
		Heartbeat:            30,
		ConnectionTimeout:    10,
		SleepOnErrorInterval: 100,
		MaxConnectionCount:   1,
		MaxCacheChannelCount: 5,
	}

	var cp = &tcr.ConnectionPool{}
	var err error
	for i := 0; i < 5; i++ {
		cp, err = tcr.NewConnectionPool(config)
		if err != nil {
			if i == 4 {
				failOnError(err, "Couldn't connect to RabbitMQ")
			}
			log.Println("Couldn't connect, retrying in 5 secs")
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
	return cp
}

// NewDeviceRPC calls AMQP to register new device
func NewDeviceRPC(cp *tcr.ConnectionPool, jsonObj newRegisterJSON) (res string, code int) {
	// Getting channel from Connection Pool
	chanHost := cp.GetChannelFromPool()

	// Declaring consume queue for RPC callback
	q, err := chanHost.Channel.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	// Consuming on that queue waiting for response
	msgs, err := chanHost.Channel.Consume(
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
	failOnError(err, "Couldn't marshal JSON")

	err = chanHost.Channel.Publish(
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

			// Process header
			switch codeU := headers["status"].(type) {
			case int:
				code = codeU
			case int8:
				code = int(codeU)
			case uint8:
				code = int(codeU)
			case int16:
				code = int(codeU)
			default:
				code = 0
			}
			break
		}
	}

	// We delete the callback queue
	_, err = chanHost.Channel.QueueDelete(q.Name, true, true, false)
	if err == nil {
		log.Println("Error removing unused queue")
	}

	// We have to return channel to the pool
	cp.ReturnChannel(chanHost, err != nil)

	return
}

// NewReading calls the HistoricRegistry Service to save a new reading
func NewReading(cp *tcr.ConnectionPool, readings readingsJSON) error {
	// Getting a channel from the pool
	chanHost := cp.GetChannelFromPool()

	// Declaring the exchange to use
	chanHost.Channel.ExchangeDeclare(
		"new_data",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	// Convert JSON object to []bytes
	jsonBytes, err := json.Marshal(readings)
	failOnError(err, "Couldn't marshal JSON")

	// Publish the message
	err = chanHost.Channel.Publish(
		"new_data",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBytes,
		},
	)
	cp.ReturnChannel(chanHost, err != nil)

	return err
}
