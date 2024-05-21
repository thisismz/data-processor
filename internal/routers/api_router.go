package routers

import "github.com/gofiber/fiber/v2"

type ApiRouter struct {
}

func (api *ApiRouter) InstallRouter(app *fiber.App) {
	//TODO: implement business logic function call
	app.Post("/data")
}
func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}
