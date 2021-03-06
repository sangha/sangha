package main

import (
	"os"

	"gitlab.techcultivation.org/sangha/mq"
	"gitlab.techcultivation.org/sangha/sangha/config"
	"gitlab.techcultivation.org/sangha/sangha/db"
	"gitlab.techcultivation.org/sangha/sangha/logger"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configFile, logLevelStr string

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
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func initConfig() {
	logLevel, err := log.ParseLevel(logLevelStr)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)

	config.ParseSettings(configFile)
	logger.SetupLogger(config.Settings.Connections.Logger.Protocol,
		config.Settings.Connections.Logger.Address,
		"sangha")

	log.Infoln("Starting sangha JSON API")

	db.SetupPostgres(config.Settings.Connections.PostgreSQL)
	mq.SetupAMQP(config.Settings.Connections.AMQP)
	db.GetDatabase()
	db.InitDatabase()
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.json", "use this config file (JSON format)")
	RootCmd.PersistentFlags().StringVarP(&logLevelStr, "loglevel", "l", "info", "log level")
}
