package main

import (
	"fmt"
	"net/http"
	"code.google.com/p/gorest"
	"code.google.com/p/goprotobuf/proto"
	"github.com/piffio/owm/protobuf"
)

var outChan chan []byte

func (serv OwmService) PostResults(testResults protobuf.TestResults) {
	message := &protobuf.TestResultsProto {
		AgentId: proto.Uint64(testResults.AgentId),
		URI: proto.String(testResults.URI),
		Timestamp: proto.String(testResults.Timestamp),
		TestData: proto.String(testResults.TestData),
	}

	data, err := proto.Marshal(message)

	if err != nil {
		LogErr("%s", fmt.Errorf("Can't Marshall message: %s", err))
		return
	}

	outChan <- data

	serv.ResponseBuilder().SetResponseCode(200)
	return
}

type OwmService struct {
	gorest.RestService `root:"/owm/" consumes:"application/json" produces:"application/json"`

	postResults gorest.EndPoint `method:"POST" path:"/postResults/" postdata:"TestResults"`
}

func ListenerWorker(cfg *Configuration, listenerStatus chan string, amqpMessages chan []byte) {
	LogDbg("Initializing Listener Worker")

	outChan = amqpMessages

	gorest.RegisterService(new(OwmService))
	http.Handle("/",gorest.Handle())
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Listener.Port), nil)
}
