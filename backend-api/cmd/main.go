package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/zeusito/narvi/internal/healthcheck"
	"github.com/zeusito/narvi/pkg/configurer"
	"github.com/zeusito/narvi/pkg/logger"
	"github.com/zeusito/narvi/pkg/router"
)

func main() {
	// Parse flags
	cfgPath := flag.String("config", "resources/config.toml", "Path to the configuration file")
	flag.Parse()

	// Setup logger
	logger.MustConfigure()

	// Load config
	configStore, err := configurer.LoadConfigurations(*cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading configurations")
	}

	// Init Http router
	myRouter := router.NewHTTPRouter(configStore.Server)

	// Init Modules Here
	_ = healthcheck.NewModule(myRouter.Mux)

	// Start server in background
	go myRouter.Start()

	// Graceful shutdown
	gracefulShutdown(myRouter)

}

func gracefulShutdown(myRouter *router.HTTPRouter) {
	// Wait for the interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// Signal acquired, starting to shut down all systems
	log.Warn().Msg("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	myRouter.Shutdown(ctx)
}
