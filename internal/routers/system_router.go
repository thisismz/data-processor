package routers

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

type SystemRouter struct {
}

func (h SystemRouter) InstallRouter(app *fiber.App) {
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	app.Get("docs/swagger/*", swagger.HandlerDefault)
}
func NewSystemRouter() *SystemRouter {
	return &SystemRouter{}
}
