package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/gofiber/fiber/v2"
	"github.com/thisismz/data-processor/pkg/databases"
	"github.com/thisismz/data-processor/pkg/env"
)

func gracefullyShutdown(app *fiber.App) {
	// Listen from a different goroutine
	go func() {
		if err := app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", ""), env.GetEnv("APP_PORT", "80"))); err != nil {
			log.Err(err).Msg("Shutting down")
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	log.Info().Msg("Gracefully shutting down...")
	databases.CloseDatabase()
	_ = app.Shutdown()

}
