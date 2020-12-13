package server

import (
	"github.com/gofiber/fiber/v2"
)

// APIResponse is a http api response JSON structure
type APIResponse struct {
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
}

// APIInternalServerError returns the fiber function for
// internal server error including the JSON data
func APIInternalServerError(c *fiber.Ctx) error {
	return c.Status(500).JSON(APIResponse{
		Status: "ERR",
		Msg:    "Internal server error.",
		Data:   nil,
	})
}
