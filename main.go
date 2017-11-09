package main

import (
	"fmt"
	"os"

	"gitlab.techcultivation.org/sangha/sangha/config"
	"gitlab.techcultivation.org/sangha/sangha/db"
	"gitlab.techcultivation.org/sangha/sangha/logger"

	log "github.com/sirupsen/logrus"
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
	logger.SetupLogger(config.Settings.Connections.Logger.Protocol,
		config.Settings.Connections.Logger.Address,
		"sangha")

	log.Infoln("Starting sangha")

	db.SetupPostgres(config.Settings.Connections.PostgreSQLConnection)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
