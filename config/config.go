package config

import (
	"encoding/json" // For config file
	"fmt"
	"os"
)

// Listener Configuration
type listener struct {
	Port int
}

// RabbitMQ Configuration
type amqpCfg struct {
	Type         string
	Host         string
	Port         int
	Vhost        string
	User         string
	Passwd       string
	Exchange     string
	ExchangeType string
	Queue        string
	RoutingKey   string
	BindingKey   string
	Workers      int
}

type logCfg struct {
	Type        string
	SyslogIdent string
	LogFile     string
}

// Main configuration structure
type Configuration struct {
	Debug    bool
	Listener listener
	Amqp     amqpCfg
	Log      logCfg
}

func ReadConfig(f string) *Configuration {
	file, _ := os.Open(f)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Printf("Could not open config file %s: %s\n", f, err)
		os.Exit(1)
	}
	return &configuration
}
