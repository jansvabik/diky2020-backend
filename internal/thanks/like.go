package thanks

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// addLike adds 1 like to the specified thanks
func addLike(id primitive.ObjectID) (*Thanks, error) {
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
		"$inc": bson.M{
			"likes": 1,
		},
	}, &fopts).Decode(&result)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &result, nil
}

// LikeHandler handles POST /thanks/:id/likes endpoint
func LikeHandler(c *fiber.Ctx) error {
	// create an object id from the id param string
	oid, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return server.APIError(c, "Dokument s tímto ID neexistuje.", 404)
	}

	// update the data in database
	doc, err := addLike(oid)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return server.APIError(c, "Dokument s tímto ID neexistuje.", 404)
		}
		return server.APIInternalServerError(c)
	}
	return server.APIOK(c, "Požadavek byl úspěšně zpracován.", doc)
}
