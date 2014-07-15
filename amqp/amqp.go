package amqp

import (
	"fmt"
	"os"
	"time"
	"github.com/streadway/amqp"
	"github.com/piffio/owm/log"
)

type queueHandler func(<-chan amqp.Delivery)

func publishMsg(connection *amqp.Connection, exchange string, exchangeType string, routingKey string, msg []byte) error {
	// Get a Channel
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	// Declare the Exchange
	log.LogDbg("Declaring Exchange \"%s\"", exchange)

	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
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
		exchange,     // name
		routingKey,   // routing to 0 or more queues
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

func ConsumeQueue(connection *amqp.Connection, exchange string, exchangeType string, bindingKey string, queueName string, tag string, handler queueHandler) {
	log.LogInf("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		log.LogErr("Channel: %s", err)
		os.Exit(2)
	}

	log.LogInf("got Channel, declaring Exchange (%q)", exchange)
	if err = channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		log.LogErr("Exchange Declare: %s", err)
	}

	log.LogInf("declared Exchange, declaring Queue %q", queueName)
	queue, err := channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		log.LogErr("Queue Declare: %s", err)
	}

	log.LogInf("declared Queue (%q %d messages, %d consumers), binding to Exchange (bindingKey %q)",
		queue.Name, queue.Messages, queue.Consumers, bindingKey)

	if err = channel.QueueBind(
		queue.Name, // name of the queue
		bindingKey,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		log.LogErr("Queue Bind: %s", err)
	}

	log.LogInf("Queue bound to Exchange, starting Consume (consumer tag %q)", tag)
	deliveries, err := channel.Consume(
		queue.Name, // name
		tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		log.LogErr("Queue Consume: %s", err)
	}

	go handler(deliveries)
}


func OpenConnection(amqpURI string, workerId string) (*amqp.Connection, error) {
	log.LogDbg("[%s] Connecting to %q", workerId, amqpURI)
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		log.LogErr("%s", fmt.Errorf("[%s] Connection error: %s", workerId, err))
	}

	go func() {
		err := <-connection.NotifyClose(make(chan *amqp.Error))
		log.LogErr("closing: %s", err)
		Reconnect(amqpURI, workerId, connection, err)
	}()

	return connection, err
}

func Reconnect(amqpURI string, workerId string, conn *amqp.Connection, err *amqp.Error) {
	var connErr error

	log.LogErr("Connection closed with error: [%d] %s", err.Code, err.Reason)

	for {
		conn, connErr = OpenConnection(amqpURI, workerId)
		if connErr != nil {
			log.LogWarn("Reconnect faled: [%s], sleeping...", connErr)
			time.Sleep(2)
		}
	}
}

func CleanupConnection(workerId string, connection *amqp.Connection) {
	log.LogDbg("[%s] Closing connection", workerId)
	connection.Close()
}
