package main

import (
	"code.google.com/p/getopt"
	"fmt"
	"os"
	"github.com/piffio/owmapi/owm"
)

func main() {
	// Parse cmdline args
	cmd_help := getopt.BoolLong("help", 'h', "", "Show command help")
	cmd_cfg := getopt.StringLong("config", 'c', "owmapi.json", "Config file")

	var opts = getopt.CommandLine
	opts.Parse(os.Args)

	if *cmd_help {
		getopt.Usage()
		os.Exit(0)
	}

	cfg, err := owm.ReadConfig(*cmd_cfg)

	if err != nil {
		fmt.Println("Failed to read config file", *cmd_cfg)
		os.Exit(1)
	}

	go owm.LoggerWorker(cfg)

	owm.LogDbg("Successfully parsed file %s", *cmd_cfg)
	// XXX pretty print the parsed conf in debug LogDbg(cfg)

	// Create RabbitMQ workers
	amqpStatus := make(chan int)
	amqpMessages := make(chan []byte)
	owm.LogDbg("Initializing AQMP Workers")

	i := 0
	for i < cfg.Amqp.Workers {
		go AmqpWorker(cfg, i, amqpStatus, amqpMessages)
		i++
	}

	// Wait for the initialization to be completed
	i = 0
	for i < cfg.Amqp.Workers {
		ready := <-amqpStatus
		if ready < 0 {
			owm.LogErr("Could not initialize AMQP workers, exiting")
			os.Exit(1)
		}
		i++
	}

	// Create Listener worker
	listenerStatus := make(chan string)
	go ListenerWorker(cfg, listenerStatus, amqpMessages)

	<-listenerStatus
}
