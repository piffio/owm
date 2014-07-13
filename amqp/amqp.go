package amqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/piffio/owm/config"
	"github.com/piffio/owm/log"
)

func publishMsg(cfg *config.Configuration, connection *amqp.Connection, msg []byte) error {
	// Get a Channel
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	// Declare the Exchange
	log.LogDbg("Declaring Exchange \"%s\"", cfg.Amqp.Exchange)

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

	log.LogDbg("Enable publishing confirm")
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
		log.LogDbg("confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		log.LogDbg("failed delivery of delivery tag: %d", tag)
	}

	return nil
}

func cleanupConnection(cfg *config.Configuration, workNum int, connection *amqp.Connection) {
	log.LogDbg("[Worker %d] Closing connection", workNum)
	connection.Close()
}
