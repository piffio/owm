package main

import (
	"fmt"
	"os"
	"code.google.com/p/getopt"
)

func main() {
	// Parse cmdline args
	cmd_help := getopt.BoolLong("help", 'h', "", "Show command help")
	cmd_cfg := getopt.StringLong("config", 'c', "owm-srv.json", "Config file")

	var opts = getopt.CommandLine
	opts.Parse(os.Args)

	if *cmd_help {
		getopt.Usage()
		os.Exit(0)
	}
	fmt.Println(*cmd_cfg)

	cfg, err := ReadConfig(*cmd_cfg)

	if err != nil {
		fmt.Println("Failed to read config file", *cmd_cfg)
		os.Exit(1)
	}

	if cfg.Debug {
		fmt.Println("Successfully parsed file", *cmd_cfg)
		fmt.Println(cfg)
	}

	// Create RabbitMQ workers
	aqmpStatus := make(chan int)
	aqmpMessages := make(chan string)
	if cfg.Debug {
		fmt.Println("Initializing AQMP Workers")
	}

	i := 0
	for i < cfg.Aqmp.Workers {
		go AqmpWorker(cfg, i, aqmpStatus, aqmpMessages)
		i++
	}

	// Wait for the initialization to be completed
	i = 0
	for i < cfg.Aqmp.Workers {
		ready := <-aqmpStatus
		if ready < 0 {
			fmt.Println("Could not initialize AMQP workers, exiting")
			os.Exit(1)
		}
		i++
	}

	// Create Listener worker
	listenerStatus := make(chan string)
	go ListenerWorker(cfg, listenerStatus, aqmpMessages)

	<-listenerStatus
}
