package main

import (
	"fmt"
	//"github.com/streadway/amqp"
)

func AqmpWorker (cfg *Configuration, i int, aqmpStatus chan string, aqmpMessages chan string) {
	if cfg.Debug {
		fmt.Println("Initializing AQMP Worker", i)
	}

	aqmpStatus <- "Worker initialization done"

	// Listen for new incoming messages
	for {
		message := <-aqmpMessages
		if cfg.Debug {
			fmt.Println(fmt.Sprintf("[Worker %d] Got message \"%s\"", i, message))
		}
	}
}
