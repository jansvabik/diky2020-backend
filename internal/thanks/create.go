package thanks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/db"
	"github.com/noltio/diky2020-backend/pkg/randomstring"
	"github.com/noltio/diky2020-backend/pkg/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// createPayload is a POST /thanks payload data structure
type createPayload struct {
	Name      string `json:"name" form:"name"`
	Addressee string `json:"addressee" form:"addressee"`
	Text      string `json:"text" form:"text"`
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

	// increment the total number of thanks for homepage analytics
	u := true
	opts := options.UpdateOptions{
		Upsert: &u,
	}
	_, err = db.Collection("welcome").UpdateOne(ctx, bson.M{}, bson.M{
		"$inc": bson.M{
			"count":  1,
			"amount": 0,
		},
	}, &opts)
	if err != nil {
		fmt.Println(err.Error())
	}

	return t, nil
}

// CreateHandler handles create requests for thanks
func CreateHandler(c *fiber.Ctx) error {
	// create new thanks struct
	var pl createPayload
	err := c.BodyParser(&pl)
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

	// create database document
	ts := time.Now()
	t := Thanks{
		Name:      pl.Name,
		ShortID:   randomstring.Generate(8),
		Text:      pl.Text,
		Addressee: pl.Addressee,
		Time:      &ts,
	}

	// add image url if uploaded
	if c.Locals("uploadedFile").(bool) {
		t.Image = "https://cdn.diky2020.cz/" + c.Locals("uploadedFileName").(string)
	}

	// save
	thanks, err := Create(t)
	if err != nil {
		return server.APIInternalServerError(c)
	}

	// encode to json and write response
	return server.APIOK(c, "Vaše poděkování jsme uložili, díky!", thanks)
}
