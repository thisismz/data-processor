package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thisismz/data-processor/internal/api/response"
)

func DataHandler(c *fiber.Ctx) error {
	// get x-user-id header
	userId := c.Get("X-User-ID")
	if userId == "" {
		// return error
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseHTTP{
			ErrorCode: 400,
			Message:   "x-user-id header not set",
		})
	}
	//TODO : add business logic here
	// return success
	return c.Status(fiber.StatusOK).JSON(response.ResponseHTTP{
		ErrorCode: 200,
		Message:   "success",
	})
}
