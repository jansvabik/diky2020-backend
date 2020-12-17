package recaptcha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/app"
	"github.com/noltio/diky2020-backend/pkg/server"
)

// request is a structure for extracting g-recaptcha-response field from request
type request struct {
	GRecaptchaResponse string `json:"g-recaptcha-response" form:"g-recaptcha-response"`
}

// googleResponse is a structure with Google validation response data
type googleResponse struct {
	Success    bool     `json:"success"`
	Hostname   string   `json:"hostname"`
	ErrorCodes []string `json:"error-codes"`
}

// Middleware is a fiber middleware for validating recaptchas
func Middleware(c *fiber.Ctx) error {
	// g-recaptcha-response field extraction and testing
	req := request{}
	c.BodyParser(&req)
	if req.GRecaptchaResponse == "" {
		return server.APIError(c, "Věříme vám, že nejste robot, ale i tak to prosím potvrďte.", 400)
	}

	// google validation request data
	postURL := "https://www.google.com/recaptcha/api/siteverify"
	postStr := url.Values{
		"secret":   {app.State.Cfg.Security.Recaptcha.SecretKey},
		"response": {req.GRecaptchaResponse},
		"remoteip": {c.IP()},
	}

	// validity check
	responsePost, err := http.PostForm(postURL, postStr)
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}
	defer responsePost.Body.Close()
	body, err := ioutil.ReadAll(responsePost.Body)
	if err != nil {
		fmt.Println(err.Error())
		return server.APIInternalServerError(c)
	}

	// unmarshal the response and test the success
	gres := googleResponse{}
	json.Unmarshal(body, &gres)
	if !gres.Success {
		return server.APIError(c, "Věříme vám, že nejste robot, ale i tak to prosím potvrďte.", 400)
	}

	// execute the next method in router
	return c.Next()
}
