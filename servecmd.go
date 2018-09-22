package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	swagger "github.com/emicklei/go-restful-swagger12"
	"github.com/muesli/smolder"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gitlab.techcultivation.org/sangha/sangha/config"
	"gitlab.techcultivation.org/sangha/sangha/db"
	"gitlab.techcultivation.org/sangha/sangha/resources/budgets"
	"gitlab.techcultivation.org/sangha/sangha/resources/codes"
	"gitlab.techcultivation.org/sangha/sangha/resources/payments"
	"gitlab.techcultivation.org/sangha/sangha/resources/projects"
	"gitlab.techcultivation.org/sangha/sangha/resources/searches"
	"gitlab.techcultivation.org/sangha/sangha/resources/sessions"
	"gitlab.techcultivation.org/sangha/sangha/resources/statistics"
	"gitlab.techcultivation.org/sangha/sangha/resources/transactions"
	"gitlab.techcultivation.org/sangha/sangha/resources/users"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "start the API and serve",
		Long:  `The serve command starts the API`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeServe()
		},
	}
)

func executeServe() error {
	ch, shutdownGracefully := handleSignals()

	context := &db.APIContext{
		Config: *config.Settings,
	}

	// Setup web-service
	smolderConfig := smolder.APIConfig{
		BaseURL:    config.Settings.API.BaseURL,
		PathPrefix: config.Settings.API.PathPrefix,
	}

	wsContainer := smolder.NewSmolderContainer(smolderConfig, &shutdownGracefully, ch)
	func(resources ...smolder.APIResource) {
		for _, r := range resources {
			r.Register(wsContainer, smolderConfig, context)
		}
	}(
		&sessions.SessionResource{},
		&users.UserResource{},
		&projects.ProjectResource{},
		&budgets.BudgetResource{},
		&codes.CodeResource{},
		&transactions.TransactionResource{},
		&payments.PaymentResource{},
		&statistics.StatisticsResource{},
		&searches.SearchesResource{},
	)

	if config.Settings.API.SwaggerFilePath != "" {
		wsConfig := swagger.Config{
			WebServices:     wsContainer.RegisteredWebServices(),
			WebServicesUrl:  config.Settings.API.BaseURL,
			ApiPath:         config.Settings.API.SwaggerAPIPath,
			SwaggerPath:     config.Settings.API.SwaggerPath,
			SwaggerFilePath: config.Settings.API.SwaggerFilePath,
		}
		swagger.RegisterSwaggerService(wsConfig, wsContainer)
	}

	// GlobalLog("Starting web-api...")
	server := &http.Server{Addr: config.Settings.API.Bind, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
	return nil
}

func handleSignals() (chan int, bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	shutdownGracefully := false
	requestIncChan := make(chan int)

	go func() {
		boldGreen := string(byte(27)) + "[1;32m"
		boldRed := string(byte(27)) + "[1;31m"
		boldEnd := string(byte(27)) + "[0m"

		pendingRequests := 0
		for {
			select {
			case sig := <-sigChan:
				if !shutdownGracefully {
					shutdownGracefully = true
					fmt.Printf(boldGreen+"\nGot %s signal, shutting down gracefully. Press Ctrl-C again to stop now.\n\n"+boldEnd, sig.String())
					if pendingRequests == 0 {
						os.Exit(0)
					}
				} else {
					fmt.Printf(boldRed+"\nGot %s signal, shutting down now!\n\n"+boldEnd, sig.String())
					os.Exit(0)
				}

			case inc := <-requestIncChan:
				pendingRequests += inc
				if shutdownGracefully {
					log.Infoln("Pending requests:", pendingRequests)
					if pendingRequests == 0 {
						os.Exit(0)
					}
				}
			}
		}
	}()

	return requestIncChan, shutdownGracefully
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
