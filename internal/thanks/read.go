package thanks

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/server"
	"go.mongodb.org/mongo-driver/bson"
)

// Read reads the specified amount of thanks from database
// and returns them in a slice of Thank structures
func Read(page uint, perPage uint) ([]Thanks, error) {
	// create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get data from database
	cursor, err := collection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// create an array of recipes and return them
	thanksList := []Thanks{}
	for cursor.Next(ctx) {
		var thanks Thanks
		cursor.Decode(&thanks)
		thanksList = append(thanksList, thanks)
	}
	return thanksList, nil
}

// ReadHandler handles read requests for thanks
func ReadHandler(c *fiber.Ctx) error {
	list, err := Read(1, 10)
	if err != nil {
		return server.APIInternalServerError(c)
	}
	return server.APIOK(c, "Požadavek byl úspěšně zpracován.", list)
}
