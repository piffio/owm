package amqp

import (
	"fmt"
	"github.com/streadway/amqp"
	"code.google.com/p/goprotobuf/proto"
	"github.com/piffio/owm/protobuf"
	"github.com/piffio/owm/config"
	"github.com/piffio/owm/log"
)

func AmqpPublisher(cfg *config.Configuration, workId int, amqpStatus chan int, amqpMessages chan []byte) {
	log.LogDbg("Initializing AQMP Publisher %d", workId)

	// Set up Worker connections
	// "amqp://guest:guest@localhost:5672/"
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.Amqp.User,
		cfg.Amqp.Passwd,
		cfg.Amqp.Host,
		cfg.Amqp.Port,
		cfg.Amqp.Vhost)

	log.LogDbg("[Worker %d] Connecting to %q", workId, uri)

	// XXX Move this in a seperate function to be called
	// On reconnection as well
	connection, err := amqp.Dial(uri)
	if err != nil {
		log.LogErr("%s", fmt.Errorf("[Worker %d] Connection error: %s", workId, err))
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
				log.LogErr("Could not Unmarshal message")
			}

			log.LogDbg("[Worker %d] Got message \"%+v\"", workId, data)
		}
		publishMsg(cfg, connection, message)
	}
}
