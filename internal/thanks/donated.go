package thanks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/app"
	"github.com/noltio/diky2020-backend/pkg/db"
	"github.com/noltio/diky2020-backend/pkg/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// donatedPayload is a PATCH /thanks/:id/donated payload data structure
type donatedPayload struct {
	ValidationToken string `json:"validationToken"`
	Target          string `json:"target"`
	Amount          int    `json:"amount"`
}

// patchDonated updates the document in database (it saves the donation
// type and the donated amount of the person)
func patchDonated(id primitive.ObjectID, target string, amount int) (*Thanks, error) {
	// create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// query options
	after := options.After
	fopts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	// try to update the database document
	var result Thanks
	err := collection().FindOneAndUpdate(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": bson.M{
			"donation": bson.M{
				"target": target,
				"amount": amount,
			},
		},
	}, &fopts).Decode(&result)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	// increment the total value of donations for homepage analytics
	u := true
	uopts := options.UpdateOptions{
		Upsert: &u,
	}
	_, err = db.Collection("welcome").UpdateOne(ctx, bson.M{}, bson.M{
		"$inc": bson.M{
			"count":  0,
			"amount": amount,
		},
	}, &uopts)
	if err != nil {
		fmt.Println(err.Error())
	}

	return &result, nil
}

// DonatedHandler handles read requests for thanks
func DonatedHandler(c *fiber.Ctx) error {
	// test json validity
	if !json.Valid(c.Body()) {
		return server.APIError(c, "Your JSON is not valid", 400)
	}

	// create new thanks struct
	var pl donatedPayload
	err := json.Unmarshal(c.Body(), &pl)
	if err != nil {
		return server.APIInternalServerError(c)
	}

	// trim whitespaces
	pl.ValidationToken = strings.TrimSpace(pl.ValidationToken)
	pl.Target = strings.TrimSpace(pl.Target)

	// validation
	if pl.ValidationToken != app.State.Cfg.Security.DonioValidationToken {
		return server.APIError(c, "Zadaný ověřovací token není validní.", 400)
	}

	// test allowed donation IDs
	if pl.Target != "fd938b2b-2fd3-4c93-bedb-df28ed75dc61" && pl.Target != "ab8da340-31df-4746-be22-a1faecc7d252" {
		return server.APIError(c, "ID, které jsme obdrželi, patří sbírce, kterou na své straně nepodporujeme. Požadavek byl úspěšně zpracován, data o přispění ale nebudou na straně Díky 2020 uložena.", 200)
	}

	// replace the target id by human-readable text
	pl.Target = map[string]string{
		"fd938b2b-2fd3-4c93-bedb-df28ed75dc61": "seniorům",
		"ab8da340-31df-4746-be22-a1faecc7d252": "samoživitelům",
	}[pl.Target]

	// create an object id from the id param string
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return server.APIError(c, "Dokument s tímto ID neexistuje.", 404)
	}

	// update the data in database
	list, err := patchDonated(oid, pl.Target, pl.Amount)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return server.APIError(c, "Dokument s tímto ID neexistuje.", 404)
		}
		return server.APIInternalServerError(c)
	}
	return server.APIOK(c, "Požadavek byl úspěšně zpracován.", list)
}
