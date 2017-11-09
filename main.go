package main

import (
	"fmt"
	"log/syslog"
	"os"

	"gitlab.techcultivation.org/sangha/sangha/config"
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/spf13/cobra"
)

var (
	// RootCmd is the core command used for cli-arg parsing
	RootCmd = &cobra.Command{
		Use:   "sangha",
		Short: "sangha JSON API server",
		Long: "sangha is the JSON API server of the sangha project\n" +
			"Complete documentation is available at https://gitlab.techcultivation.org/sangha/sangha",
		SilenceErrors: false,
		SilenceUsage:  true,
	}
)

func main() {
	config.ParseSettings()

	log := logrus.New()
	hook, err := lSyslog.NewSyslogHook("tcp", "10.0.3.216:5514", syslog.LOG_INFO, "sangha")
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
	} else {
		log.Hooks.Add(hook)
	}

	db.SetupPostgres(config.Settings.Connections.PostgreSQLConnection)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
