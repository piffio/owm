package main

import (
	"github.com/piffio/owm/amqp"
	"github.com/piffio/owm/config"
	"github.com/piffio/owm/log"
	"github.com/piffio/owm/protobuf"
	"github.com/rcrowley/go-tigertonic/mocking"
	"net/http"
	"testing"
)

var (
	cfg          *config.Configuration
	amqpStatus   = make(chan int)
	amqpMessages = make(chan []byte)
	err          error
)

func TestReadConfig(t *testing.T) {
	cfg_file := "../conf/owmapi.json"
	cfg = config.ReadConfig(cfg_file)

	if cfg.Amqp.Workers <= 0 {
		t.Fatal("Missing AMQP configuration")
	}

	go log.LoggerWorker(cfg)
}

func Test2AMQPUp(t *testing.T) {
	i := 0
	for i < cfg.Amqp.Workers {
		go amqp.AmqpPublisher(cfg, i, amqpStatus, amqpMessages)
		i++
	}

	i = 0
	for i < cfg.Amqp.Workers {
		ready := <-amqpStatus
		if ready != 1 {
			t.Fatal(ready)
		}
		i++
	}
}

func TestListernerUp(t *testing.T) {
	listenerStatus := make(chan string)
	go ListenerWorker(cfg, listenerStatus, amqpMessages)
	listenerStatus <- "status"
	status := <-listenerStatus
	if status != "running" {
		t.Fatal(status)
	}
}

func DISABLEDTestPostResults(t *testing.T) {
	code, header, response, err := postResultsHandler(
		mocking.URL(mux, "POST", "http://localhost:9999/owm/postResults"),
		mocking.Header(nil),
		&protobuf.TestResults{12345, "http://www.piffio.org", "13243433", "somethingsomething"},
	)
	if nil != err {
		t.Fatal(err)
	}
	if http.StatusCreated != code {
		t.Fatal(code)
	}
	if "http://localhost:9999/1.0/postResults/12345" != header.Get("Content-Location") {
		t.Fatal(header)
	}
	if nil == response {
		t.Fatal("Empty response!")
	}
	if 12345 != response.AgentId || "http://www.piffio.org" != response.URI {
		t.Fatal(response)
	}
}
