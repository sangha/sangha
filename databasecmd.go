package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.techcultivation.org/sangha/sangha/db"
)

var (
	databaseCmd = &cobra.Command{
		Use:   "database",
		Short: "manage database",
		Long:  `The database command is used to init or migrate the database`,
		RunE:  nil,
	}
	databaseInitCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize the database",
		Long:  `The init command initializes the database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDatabaseInit()
		},
	}
	databaseWipeCmd = &cobra.Command{
		Use:   "wipe",
		Short: "wipe the database",
		Long:  `The wipe command wipes the entire database and drops all tables`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDatabaseWipe()
		},
	}
)

func init() {
	databaseCmd.AddCommand(databaseInitCmd)
	databaseCmd.AddCommand(databaseWipeCmd)
	RootCmd.AddCommand(databaseCmd)
}

func executeDatabaseInit() error {
	log.Println("Init database")

	db.GetDatabase()
	db.InitDatabase()

	return nil
}

func executeDatabaseWipe() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you really want to wipe the entire database?\nEnter 'SELFDESTRUCT' to confirm: ")
	text, _ := reader.ReadString('\n')

	if strings.TrimSpace(text) != "SELFDESTRUCT" {
		return errors.New("Wiping database requires user confirmation")
	}

	log.Println("Wiping database")

	db.GetDatabase()
	db.WipeDatabase()

	return nil
}
