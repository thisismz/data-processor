package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thisismz/data-processor/internal/api"
)

type ApiRouter struct {
}

func (a *ApiRouter) InstallRouter(app *fiber.App) {
	app.Post("/data", api.DataHandler)
}
func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}
