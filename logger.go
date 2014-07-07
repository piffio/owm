package main

import (
	"os"
	"fmt"
	"log"
	"log/syslog"
)

var logDbgChan = make(chan string)
var logInfChan = make(chan string)
var logWarnChan = make(chan string)
var logErrChan = make(chan string)

func LogDbg(format string,  params...interface{}) {
	if len(params) > 0 {
		logDbgChan <-fmt.Sprintf(format, params[0])
	} else {
		logDbgChan <-format
	}
}

func LogInf(format string,  params...interface{}) {
	if len(params) > 0 {
		logInfChan <-fmt.Sprintf(format, params[0])
	} else {
		logInfChan <-format
	}
}

func LogWarn(format string,  params...interface{}) {
	if len(params) > 0 {
		logWarnChan <-fmt.Sprintf(format, params[0])
	} else {
		logWarnChan <-format
	}
}

func LogErr(format string,  params...interface{}) {
	if len(params) > 0 {
		logErrChan <-fmt.Sprintf(format, params[0])
	} else {
		logErrChan <-format
	}
}

// Main goroutine that handles log spooling
func LoggerWorker(cfg *Configuration) {
	var logf *log.Logger
	var logs *syslog.Writer
	var logPrefix string

	// In debug mode we always send the same
	// message to a log file
	if cfg.Log.Type == "file" || cfg.Debug {
		if cfg.Log.LogFile != "" {
			if cfg.Log.SyslogIdent != "" {
				logPrefix = cfg.Log.SyslogIdent
			} else {
				logPrefix = "owmlogger"
			}

			file, err := os.Create(cfg.Log.LogFile)
			if err != nil { panic(err) }

			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()

			logf = log.New(file, logPrefix, log.Lshortfile)
		}
	}

	if cfg.Log.Type == "syslog" {
		var err error
		logs, err = syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, cfg.Log.SyslogIdent)
		if err != nil {
			fmt.Println("Cannot initialize Syslog support")
		}
	}

	for {
		select {
			case msg := <-logInfChan:
			if logf != nil {
				logf.Println(msg)
			}
			if logs != nil {
				logs.Info(msg)
			}
			case msg := <-logDbgChan:
			if cfg.Debug {
				if logf != nil {
					logf.Println(msg)
				}
				if logs != nil {
					logs.Debug(msg)
				}
			}
		}
	}
}
