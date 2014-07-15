package amqp

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/piffio/owm/config"
	"github.com/piffio/owm/log"
	"github.com/piffio/owm/protobuf"
	"github.com/streadway/amqp"
	"os"
)

func AmqpWorker(cfg *config.Configuration, workId int, amqpStatus chan int) {
	var workerId = fmt.Sprintf("Worker %d", workId)
	log.LogDbg("Initializing AQMP %s", workerId)

	// Set up Worker connections
	// "amqp://guest:guest@localhost:5672/"
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.Amqp.User,
		cfg.Amqp.Passwd,
		cfg.Amqp.Host,
		cfg.Amqp.Port,
		cfg.Amqp.Vhost)

	connection, err := OpenConnection(uri, workerId)
	if err != nil {
		os.Exit(1)
	}
	log.LogDbg("[%s] Connection established", workerId)
	defer CleanupConnection(workerId, connection)

	amqpStatus <- 1

	ConsumeQueue(connection, cfg.Amqp.Exchange, cfg.Amqp.ExchangeType, cfg.Amqp.BindingKey, cfg.Amqp.Queue, workerId, handle)

	select {}
}

/* XXX do something useful with the message */
func handle(deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		message := new(protobuf.TestResultsProto)
		err := proto.Unmarshal(d.Body, message)
		if err != nil {
			log.LogErr("Could not marshal message: %q", d.Body)
		} else {
			log.LogInf(
				"got %dB delivery: [%v] \"%+v\"",
				len(d.Body),
				d.DeliveryTag,
				message,
			)
		}
		d.Ack(false)
	}
}
