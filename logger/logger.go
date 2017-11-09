package logger

import (
	"fmt"
	"log/syslog"

	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

func SetupLogger(protocol, host, app string) {
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
