package main

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/piffio/owm/config"
	"github.com/piffio/owm/log"
	"github.com/piffio/owm/protobuf"
	"github.com/rcrowley/go-tigertonic"
	"net/http"
	"net/url"
)

var (
	outChan chan []byte
	mux     *tigertonic.TrieServeMux
)

func postResultsHandler(u *url.URL, h http.Header, rq *protobuf.TestResults) (int, http.Header, *protobuf.TestResults, error) {
	message := &protobuf.TestResultsProto{
		AgentId:   proto.Uint64(rq.AgentId),
		URI:       proto.String(rq.URI),
		Timestamp: proto.String(rq.Timestamp),
		TestData:  proto.String(rq.TestData),
	}

	data, err := proto.Marshal(message)

	if err != nil {
		log.LogErr("%s", fmt.Errorf("Can't Marshall message: %s", err))
		// XXX return
	}

	outChan <- data

	return http.StatusCreated, http.Header{
		"Content-Location": {fmt.Sprintf(
			"%s://%s/1.0/postResults/%s", // TODO Don't hard-code this.
			u.Scheme,
			u.Host,
			rq.AgentId,
		)},
	}, rq, nil
}

func ListenerStatus(listenerStatus chan string) {
	select {
	case tag := <-listenerStatus:
		log.LogDbg("Received status request: %s", tag)
		listenerStatus <- "running"
	}
}

func ListenerWorker(cfg *config.Configuration, listenerStatus chan string, amqpMessages chan []byte) {
	log.LogDbg("Initializing Listener Worker")

	outChan = amqpMessages

	mux = tigertonic.NewTrieServeMux()
	mux.HandleNamespace("/owm", mux)

	mux.Handle("GET", "/version", tigertonic.Version("0.1"))

	mux.Handle(
		"POST",
		"/postResults",
		tigertonic.Timed(tigertonic.Marshaled(postResultsHandler), "POST-results", nil),
	)

	go ListenerStatus(listenerStatus)

	tigertonic.NewServer(fmt.Sprintf(":%d", cfg.Listener.Port), tigertonic.Logged(mux, nil)).ListenAndServe()
}
