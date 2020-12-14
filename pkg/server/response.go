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
	err := c.Status(500).JSON(APIResponse{
		Status: "ERR",
		Msg:    "Omlouváme se, nastala u nás nečekaná chyba. Zkuste to prosím znovu za chvíli.",
		Data:   nil,
	})
	if err != nil {
		return err
	}

	// fix old browsers (like Safari) problems with encoding
	// because of this line, we don't follow the RFC standard
	c.Set("Content-Type", fiber.MIMEApplicationJSONCharsetUTF8)
	return nil
}

// APIOK returns the fiber function for status OK response
func APIOK(c *fiber.Ctx, msg string, data interface{}) error {
	err := c.Status(200).JSON(APIResponse{
		Status: "OK",
		Msg:    msg,
		Data:   data,
	})
	if err != nil {
		return err
	}

	// fix old browsers (like Safari) problems with encoding
	// because of this line, we don't follow the RFC standard
	c.Set("Content-Type", fiber.MIMEApplicationJSONCharsetUTF8)
	return nil
}

// APIError returns the fiber function for error responses
func APIError(c *fiber.Ctx, msg string, code int) error {
	err := c.Status(code).JSON(APIResponse{
		Status: "ERR",
		Msg:    msg,
	})
	if err != nil {
		return err
	}

	// fix old browsers (like Safari) problems with encoding
	// because of this line, we don't follow the RFC standard
	c.Set("Content-Type", fiber.MIMEApplicationJSONCharsetUTF8)
	return nil
}
