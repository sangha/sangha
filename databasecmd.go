package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gosimple/slug"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.techcultivation.org/sangha/sangha/config"
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
	databaseMockCmd = &cobra.Command{
		Use:   "mock",
		Short: "generate mock-up data",
		Long:  `The mock command generates mock-up data in the database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDatabaseMock()
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
	databaseCmd.AddCommand(databaseMockCmd)
	databaseCmd.AddCommand(databaseWipeCmd)
	RootCmd.AddCommand(databaseCmd)
}

func executeDatabaseInit() error {
	log.Println("Init database")

	db.GetDatabase()
	db.InitDatabase()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("New password for user admin: ")
	password, _ := reader.ReadString('\n')

	db.GetDatabase()
	context := &db.APIContext{
		Config: *config.Settings,
	}
	ctx := context.NewAPIContext().(*db.APIContext)

	user := db.User{
		Nickname: "admin",
		Email:    "admin@techcultivation.org",
		About:    "admin",
		Address:  []string{},
		ZIP:      "",
		City:     "",
		Country:  "",
	}
	err := user.Save(ctx)
	if err != nil {
		return err
	}

	err = user.UpdatePassword(ctx, strings.TrimSpace(password))
	if err != nil {
		return err
	}

	project := db.Project{
		Slug:           "cct",
		Name:           "CCT",
		Summary:        "Center for the Cultivation of Technology",
		About:          "",
		Website:        "https://techcultivation.org",
		License:        "AGPL",
		Repository:     "https://techcultivation.org",
		Private:        false,
		PrivateBalance: true,
	}
	err = project.Save(ctx)
	if err != nil {
		return err
	}

	budget := db.Budget{
		ProjectID:      &project.ID,
		ParentID:       0,
		Name:           project.Name,
		Private:        false,
		PrivateBalance: true,
	}
	err = budget.Save(ctx)
	if err != nil {
		return err
	}
	dbudget := db.Budget{
		ProjectID:      &project.ID,
		ParentID:       0,
		Name:           "Donation Cuts",
		Private:        false,
		PrivateBalance: true,
	}
	err = dbudget.Save(ctx)
	if err != nil {
		return err
	}

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

func executeDatabaseMock() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you really want to write mock-up data to the database?\nEnter 'MOCKUP' to confirm: ")
	text, _ := reader.ReadString('\n')

	if strings.TrimSpace(text) != "MOCKUP" {
		return errors.New("Generating mock-up data requires user confirmation")
	}

	log.Println("Generating mock-up data")

	db.GetDatabase()
	context := &db.APIContext{
		Config: *config.Settings,
	}
	ctx := context.NewAPIContext().(*db.APIContext)

	gofakeit.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		code, err := mockProject(ctx)
		if err != nil {
			return err
		}

		for i := 0; i < 24; i++ {
			t, err := mockPayment(ctx, 2, i%3 == 0)
			if err != nil {
				return err
			}

			if i%2 == 0 {
				t.Pending = false
				t.Code = code
				err = t.Update(ctx)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func mockPayment(ctx *db.APIContext, budget int64, negative bool) (db.Payment, error) {
	valuei, _ := big.NewFloat(gofakeit.Price(1, 100) * 100.0).Int64()
	if negative {
		valuei *= -1
	}

	t := db.Payment{
		BudgetID:      budget,
		CreatedAt:     gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
		Amount:        valuei,
		Currency:      "EUR",
		Purpose:       gofakeit.Sentence(5),
		RemoteAccount: strconv.FormatInt(int64(gofakeit.CreditCard().Number), 10),
		RemoteBankID:  strconv.FormatInt(int64(gofakeit.CreditCard().Number), 10),
		RemoteName:    gofakeit.Name(),
		Source:        "hbci",
	}

	return t, t.Save(ctx)
}

func mockProject(ctx *db.APIContext) (string, error) {
	name := gofakeit.Company()

	project := db.Project{
		Slug:           slug.Make(name),
		Name:           name,
		Summary:        gofakeit.HipsterSentence(10),
		About:          gofakeit.HipsterParagraph(2, 4, 20, "."),
		Website:        gofakeit.URL(),
		License:        "GPL",
		Repository:     gofakeit.URL(),
		Private:        false,
		PrivateBalance: true,
	}
	err := project.Save(ctx)
	if err != nil {
		return "", err
	}

	budget := db.Budget{
		ProjectID:      &project.ID,
		ParentID:       0,
		Name:           project.Name,
		Private:        false,
		PrivateBalance: true,
	}
	err = budget.Save(ctx)
	if err != nil {
		return "", err
	}

	code, err := ctx.LoadCodeByBudgetUUID(budget.UUID)
	fmt.Println("Generated code:", code.Code)

	return code.Code, nil
}
