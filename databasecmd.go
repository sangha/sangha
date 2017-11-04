package main

import (
	"fmt"

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
)

func init() {
	databaseCmd.AddCommand(databaseInitCmd)
	RootCmd.AddCommand(databaseCmd)
}

func executeDatabaseInit() error {
	fmt.Printf("Init database\n")

	db.GetDatabase()
	db.InitDatabase()

	return nil
}
