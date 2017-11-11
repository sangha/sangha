package logger

import (
	"fmt"
	"log/syslog"

	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

// SetupLogger initializes logrus
func SetupLogger(protocol, host, app string) {
	if len(protocol) == 0 {
		return
	}

	hook, err := lSyslog.NewSyslogHook(protocol, host, syslog.LOG_INFO, app)
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
	} else {
		log.AddHook(hook)
	}
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}
