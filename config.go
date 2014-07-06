package main

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
	RoutingKey   string
	Workers      int
}

// Main configuration structure
type Configuration struct {
	Debug    bool
	Listener listener
	Amqp     amqpCfg
}

func ReadConfig(f string) (*Configuration, error) {
	file, _ := os.Open(f)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return &configuration, err
}
