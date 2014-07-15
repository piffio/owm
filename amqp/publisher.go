package amqp

import (
	"fmt"
	"os"
	"code.google.com/p/goprotobuf/proto"
	"github.com/piffio/owm/protobuf"
	"github.com/piffio/owm/config"
	"github.com/piffio/owm/log"
)

func AmqpPublisher(cfg *config.Configuration, id int, amqpStatus chan int, amqpMessages chan []byte) {
	var workerId = fmt.Sprintf("Publisher %d", id)
	log.LogDbg("Initializing AQMP %s", workerId)

	// Set up Publisher connections
	// syntax: "amqp://user:passwd@host:port/vhost"
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

	// Positive value means success
	amqpStatus <- 1

	// Listen for new incoming messages
	select {
	case message := <-amqpMessages:
		if cfg.Debug {
			data := new(protobuf.TestResultsProto)

			err := proto.Unmarshal(message, data)
			if err != nil {
				log.LogErr("Could not Unmarshal message")
			}

			log.LogDbg("[%s] Got message \"%+v\"", workerId, data)
		}
		publishMsg(connection, cfg.Amqp.Exchange, cfg.Amqp.ExchangeType, cfg.Amqp.RoutingKey, message)
	}
}
