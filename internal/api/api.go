package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thisismz/data-processor/internal/api/request"
	"github.com/thisismz/data-processor/internal/api/response"
	"github.com/thisismz/data-processor/internal/service"
)

// DataHandler is the handler for data
// @Summary Send data to the server
// @Description Send data to the server and return the result
// @Tags Data
// @Accept json
// @Produce json
// @Param data body request.DataRequest true "Data to send"
// @Success 200 {object} response.ResponseHTTP
// @Failure 400 {object} response.ResponseHTTP
// @Failure 429 {object} response.ResponseHTTP
// @Failure 500 {object} response.ResponseHTTP
// @Router /data [post]
func DataHandler(c *fiber.Ctx) error {
	var dataRequest request.DataRequest
	if err := c.BodyParser(&dataRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseHTTP{
			ErrorCode: 400,
			Message:   err.Error(),
		})
	}
	dataSize := int64(len(dataRequest.Payload))
	err := service.UserLimitsCheck(dataRequest.UserID, dataRequest.DataID, dataSize)
	if err != nil {
		return c.Status(fiber.StatusTooManyRequests).JSON(response.ResponseHTTP{
			ErrorCode: 429,
			Message:   err.Error(),
		})
	}

	err = service.SendDataToQueue(dataRequest.UserID, dataRequest.DataID, dataRequest.Payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseHTTP{
			ErrorCode: 500,
			Message:   err.Error(),
		})
	}
	// return success
	return c.Status(fiber.StatusOK).JSON(response.ResponseHTTP{
		ErrorCode: 200,
		Message:   "success",
	})
}
