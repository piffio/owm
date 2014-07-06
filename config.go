package main

import (
	"fmt"
	"os"
	"encoding/json" // For config file
)

// Listener Configuration
type listener struct {
	Port int
}

// RabbitMQ Configuration
type aqmp struct {
	Type string
	Host string
	Port int
	Vhost string
	User string
	Passwd string
	Workers int
}

// Main configuration structure
type Configuration struct {
	Debug bool
	Listener listener
	Aqmp aqmp
}

func ReadConfig (f string) (*Configuration, error) {
	file, _ := os.Open(f)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return &configuration, err
}
