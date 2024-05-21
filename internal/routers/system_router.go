package routers

import "github.com/gofiber/fiber/v2"

type SystemRouter struct {
}

func (h SystemRouter) InstallRouter(app *fiber.App) {
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
}
func NewSystemRouter() *SystemRouter {
	return &SystemRouter{}
}
