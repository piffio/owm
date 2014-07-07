package main

import (
	"fmt"
	"net/http"
	"code.google.com/p/gorest"
)

var outChan chan TestResults

/*
func handler(w http.ResponseWriter, r *http.Request, amqpMessages chan TestResults) {
	LogDbg("Inside handler")
	fmt.Fprintf(w, "Hello world from my Go program!\n")

	amqpMessages <- "got one request"
}*/

func (serv OwmService) PostResults(testResults TestResults) {
	outChan <- testResults

	serv.ResponseBuilder().SetResponseCode(200)
	return
}

type OwmService struct {
	gorest.RestService `root:"/owm/" consumes:"application/json" produces:"application/json"`

	//agentDetails gorest.EndPoint `method:"GET" path:"/agent/{Id:int}" output:"Agent"`
	postResults gorest.EndPoint `method:"POST" path:"/postResults/" postdata:"TestResults"`
}

func ListenerWorker(cfg *Configuration, listenerStatus chan string, amqpMessages chan TestResults) {
	LogDbg("Initializing Listener Worker")

	outChan = amqpMessages

	gorest.RegisterService(new(OwmService))
	http.Handle("/",gorest.Handle())
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, amqpMessages)
	})*/
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Listener.Port), nil)
}
