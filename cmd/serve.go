package cmd

import (
	"fmt"
	"github.com/s4kibs4mi/jally-commerce-bot/config"
	"github.com/s4kibs4mi/jally-commerce-bot/handlers"
	"github.com/s4kibs4mi/jally-commerce-bot/log"
	"github.com/s4kibs4mi/jally-commerce-bot/processors"
	"github.com/s4kibs4mi/jally-commerce-bot/services"
	"github.com/s4kibs4mi/jally-commerce-bot/services/messenger"
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

	var processor processors.IStateProcessor
	var r handlers.IRouter

	if config.App().ActiveProcessor == "facebook" {
		fbMessenger := &messenger.Messenger{
			AccessToken: config.App().FacebookAccessToken,
			Debug:       messenger.DebugAll,
		}
		processor, err = processors.NewFacebookStateProcessor(config.App(), stateService, shopemaaService, fbMessenger)
		if err != nil {
			panic(err)
		}
		r = handlers.NewRouter(config.App(), processor, shopemaaService, fbMessenger)
	} else if config.App().ActiveProcessor == "twilio" {
		twilioService := services.NewTwilioService(config.App())
		processor, err = processors.NewTwilioStateProcessor(config.App(), stateService, shopemaaService, twilioService)
		if err != nil {
			panic(err)
		}
		r = handlers.NewRouter(config.App(), processor, shopemaaService, nil)
	} else {
		panic("Unknown processor selected")
	}

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
