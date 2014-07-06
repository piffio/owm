package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request, amqpMessages chan string) {
	fmt.Println("Inside handler")
	fmt.Fprintf(w, "Hello world from my Go program!\n")

	amqpMessages <-"got one request"
}

func ListenerWorker(cfg *Configuration, listenerStatus chan string, amqpMessages chan string) {
	if cfg.Debug {
		fmt.Println("Initializing Listener Worker")
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, amqpMessages)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Listener.Port), nil)
}
