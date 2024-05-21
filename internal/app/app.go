package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/thisismz/data-processor/pkg/env"
)

const idleTimeout = 30 * time.Second

func Run() {
	env.SetupEnvFile()
	app := fiber.New(fiber.Config{
		IdleTimeout: idleTimeout,
		BodyLimit:   50 * 1024 * 1024,
	})
	app.Use(recover.New())
	gracefullyShutdown(app)
}
