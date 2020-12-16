package welcome

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/server"
	"go.mongodb.org/mongo-driver/bson"
)

// Get gets the homepage analytics data from database
// and returns them withing the welcomeData structure
func Get() (WData, error) {
	// create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get the data
	var result WData
	err := collection().FindOne(ctx, bson.M{}).Decode(&result)
	if err != nil {
		return WData{}, err
	}

	return result, nil
}

// ReadHandler handles read requests for welcome
func ReadHandler(c *fiber.Ctx) error {
	// get welcome data and create the response
	wdata, err := Get()
	if err != nil {
		return server.APIInternalServerError(c)
	}

	// create and send the response structure
	response := map[string]interface{}{
		"totalThanks": wdata.Count,
		"donated":     wdata.Amount,
	}
	return server.APIOK(c, "Požadavek byl úspěšně zpracován.", response)
}
