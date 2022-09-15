package cmd

import (
	"fmt"
	"github.com/s4kibs4mi/twilfe/config"
	"github.com/s4kibs4mi/twilfe/handlers"
	"github.com/s4kibs4mi/twilfe/log"
	"github.com/s4kibs4mi/twilfe/processors"
	"github.com/s4kibs4mi/twilfe/services"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve starts http server",
	Run:   serve,
}

func serve(cmd *cobra.Command, args []string) {
	stateService := services.NewStateService()
	shopemaaService, err := services.NewShopemaaService(config.App())
	if err != nil {
		panic(err)
	}

	twilioService := services.NewTwilioService(config.App())
	processorService := processors.NewStateProcessor(stateService, shopemaaService, twilioService)
	r := handlers.NewRouter(processorService, shopemaaService)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	log.Log().Infoln("Starting HTTP server...")

	go func() {
		err := r.Start(fmt.Sprintf("%s:%d", config.App().Host, config.App().Port))
		if err != nil {
			log.Log().Infoln("Failed to start HTTP server.")
			panic(err)
		}
	}()

	<-stop

	log.Log().Infoln("Stopping HTTP server...")

	r.Stop()
}
