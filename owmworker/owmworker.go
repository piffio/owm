package main

import (
	"os"
	"code.google.com/p/getopt"
	"github.com/piffio/owm/amqp"
	"github.com/piffio/owm/config"
	"github.com/piffio/owm/log"
)

func main() {
	// Parse cmdline args
	cmd_help := getopt.BoolLong("help", 'h', "", "Show command help")
	cmd_cfg := getopt.StringLong("config", 'c', "owmworker.json", "Config file")

	var opts = getopt.CommandLine
	opts.Parse(os.Args)

	if *cmd_help {
		getopt.Usage()
		os.Exit(0)
	}

	cfg := config.ReadConfig(*cmd_cfg)

	go log.LoggerWorker(cfg)

	amqpStatus := make(chan int)
	//amqpMessages := make(chan []byte)
	log.LogDbg("Initializing AQMP Workers")

	i := 0
	for i < cfg.Amqp.Workers {
		go amqp.AmqpWorker(cfg, i, amqpStatus)
		i++
	}

	i = 0
	for i < cfg.Amqp.Workers {
		ready := <-amqpStatus
		if ready < 0 {
			log.LogErr("Could not initialize AMQP workers, exiting")
			os.Exit(1)
		}
		i++
	}

	select {}
}
