package app

import (
	"os"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/thisismz/data-processor/pkg/env"
)

const idleTimeout = 30 * time.Second

func Run() {
	env.SetupEnvFile()
	// zero log config
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// fiber config
	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
		BodyLimit:   50 * 1024 * 1024,
	})

	// fiber middleware
	app.Use(recover.New())
	// fiber logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	gracefullyShutdown(app)
}
