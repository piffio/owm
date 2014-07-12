package main

import (
	"fmt"
	"net/http"
	"net/url"
	"code.google.com/p/goprotobuf/proto"
	"github.com/piffio/owm/protobuf"
	"github.com/rcrowley/go-tigertonic"
)

var outChan chan []byte

func postResultsHandler(u *url.URL, h http.Header, rq *protobuf.TestResults) (int, http.Header, *protobuf.TestResults, error) {
	message := &protobuf.TestResultsProto {
		AgentId: proto.Uint64(rq.AgentId),
		URI: proto.String(rq.URI),
		Timestamp: proto.String(rq.Timestamp),
		TestData: proto.String(rq.TestData),
	}

	data, err := proto.Marshal(message)

	if err != nil {
		LogErr("%s", fmt.Errorf("Can't Marshall message: %s", err))
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

func ListenerWorker(cfg *Configuration, listenerStatus chan string, amqpMessages chan []byte) {
	LogDbg("Initializing Listener Worker")

	outChan = amqpMessages

	mux := tigertonic.NewTrieServeMux()
	mux.HandleNamespace("/owm", mux)

	mux.Handle("GET", "/version", tigertonic.Version("0.1"))

	mux.Handle(
		"POST",
		"/postResults",
		tigertonic.Timed(tigertonic.Marshaled(postResultsHandler), "POST-results", nil),
	)

	tigertonic.NewServer(fmt.Sprintf(":%d", cfg.Listener.Port), tigertonic.Logged(mux, nil)).ListenAndServe()
}
