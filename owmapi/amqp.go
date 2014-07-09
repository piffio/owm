package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"code.google.com/p/goprotobuf/proto"
	"github.com/piffio/owmapi/protobuf"
)

func publishMsg(cfg *Configuration, connection *amqp.Connection, msg []byte) error {
	// Get a Channel
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	// Declare the Exchange
	LogDbg("Declaring Exchange \"%s\"", cfg.Amqp.Exchange)

	if err := channel.ExchangeDeclare(
		cfg.Amqp.Exchange,     // name
		cfg.Amqp.ExchangeType, // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	LogDbg("Enable publishing confirm")
	if err := channel.Confirm(false); err != nil {
		return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
	}

	ack, nack := channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	// Send the Message
	if err = channel.Publish(
		cfg.Amqp.Exchange,   // publish to an exchange
		cfg.Amqp.RoutingKey, // routing to 0 or more queues
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			Body:            msg,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	// Wait for message confirmation
	select {
	case tag := <-ack:
		LogDbg("confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		LogDbg("failed delivery of delivery tag: %d", tag)
	}

	return nil
}

func cleanupConnection(cfg *Configuration, workNum int, connection *amqp.Connection) {
	LogDbg("[Worker %d] Closing connection", workNum)
	connection.Close()
}

func AmqpWorker(cfg *Configuration, workId int, amqpStatus chan int, amqpMessages chan []byte) {
	LogDbg("Initializing AQMP Worker %d", workId)

	// Set up Worker connections
	// "amqp://guest:guest@localhost:5672/"
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.Amqp.User,
		cfg.Amqp.Passwd,
		cfg.Amqp.Host,
		cfg.Amqp.Port,
		cfg.Amqp.Vhost)

	LogDbg("[Worker %d] Connecting to %q", workId, uri)

	// XXX Move this in a seperate function to be called
	// On reconnection as well
	connection, err := amqp.Dial(uri)
	if err != nil {
		LogErr("%s", fmt.Errorf("[Worker %d] Connection error: %s", workId, err))
		amqpStatus <- -1
	}
	defer cleanupConnection(cfg, workId, connection)

	// Positive value means success
	// TODO: Use an enum to allow for different states
	amqpStatus <- 1

	// Listen for new incoming messages
	for {
		message := <-amqpMessages

		if cfg.Debug {
			data := new(protobuf.TestResultsProto)

			err := proto.Unmarshal(message, data)
			if err != nil {
				LogErr("Could not Unmarshal message")
			}

			LogDbg("[Worker %d] Got message \"%+v\"", workId, data)
		}
		publishMsg(cfg, connection, message)
	}
}
