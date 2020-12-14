package thanks

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/server"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// createPayload is a POST /thanks payload data structure
type createPayload struct {
	Name      string `json:"name"`
	Addressee string `json:"addressee"`
	Text      string `json:"text"`
}

// Create creates new thanks in database
func Create(t Thanks) (Thanks, error) {
	// create new recipe document
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := collection().InsertOne(ctx, t)
	if err != nil {
		return t, err
	}

	// add id to the recipe structure
	oid := result.InsertedID.(primitive.ObjectID)
	t.ID = oid

	return t, nil
}

// CreateHandler handles create requests for thanks
func CreateHandler(c *fiber.Ctx) error {
	// test json validity
	if !json.Valid(c.Body()) {
		return c.Status(400).JSON(server.APIResponse{
			Status: "ERR",
			Msg:    "Your JSON is not valid.",
		})
	}

	// create new thanks struct
	var pl createPayload
	err := json.Unmarshal(c.Body(), &pl)
	if err != nil {
		return server.APIInternalServerError(c)
	}

	// trim whitespaces
	pl.Name = strings.TrimSpace(pl.Name)
	pl.Text = strings.TrimSpace(pl.Text)
	pl.Addressee = strings.TrimSpace(pl.Addressee)

	// name validation
	if len(pl.Name) < 3 {
		return c.JSON(server.APIResponse{
			Status: "ERR",
			Msg:    "Jméno by mělo mít alespoň tři znaky.",
		})
	}

	// addressee validation
	if len(pl.Addressee) < 3 {
		return c.JSON(server.APIResponse{
			Status: "ERR",
			Msg:    "Příjemce poděkování by měl mít alespoň tři znaky.",
		})
	}

	// text validation
	if len(pl.Text) < 3 {
		return c.JSON(server.APIResponse{
			Status: "ERR",
			Msg:    "Zkusíte vymyslet o trošku delší poděkování, ať to stojí za to? Zkuste ho třeba někomu věnovat nebo vyjmenovat věci, za které jste vděčni. :)",
		})
	}

	// save to database
	ts := time.Now()
	thanks, err := Create(Thanks{
		Name:      pl.Name,
		Text:      pl.Text,
		Addressee: pl.Addressee,
		Time:      &ts,
	})
	if err != nil {
		return server.APIInternalServerError(c)
	}

	// encode to json and write response
	return c.JSON(server.APIResponse{
		Status: "OK",
		Msg:    "Recipe was saved successfully.",
		Data:   thanks,
	})
}
