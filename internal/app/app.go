package app

import (
	"os"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/thisismz/data-processor/internal/routers"
	"github.com/thisismz/data-processor/internal/service"
	"github.com/thisismz/data-processor/pkg/circuit_breaker"
	"github.com/thisismz/data-processor/pkg/databases"
	"github.com/thisismz/data-processor/pkg/env"
)

const idleTimeout = 30 * time.Second

func Run() {
	env.SetupEnvFile()
	// start databases
	databases.StartDatabase()
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

	if err := service.StorageServiceUp(); err != nil {
		log.Err(err).Msg("storage service failed")
		os.Exit(1)
	}
	if err := service.QueueServiceUp(); err != nil {
		log.Err(err).Msg("queue service failed")
		os.Exit(1)
	}
	// install routers
	routers.InstallRouter(app)
	// start data processor
	service.StartDataProcessor()
	isLeader := true
	if env.GetEnv("IS_LEADER", "TRUE") == "false" {
		isLeader = false
	}
	circute := circuit_breaker.New("redis", 0, time.Duration(10*time.Minute), isLeader)
	circute.Run()
	gracefullyShutdown(app)
}
