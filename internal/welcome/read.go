package welcome

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/internal/thanks"
	"github.com/noltio/diky2020-backend/pkg/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Get gets the homepage analytics data from database
// and returns them withing the welcomeData structure
func Get() (WData, error) {
	// create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// mongodb options
	opts := options.FindOneOptions{}

	// get the data
	var result WData
	err := collection().FindOne(ctx, bson.M{}, &opts).Decode(&result)
	if err != nil {
		return WData{}, err
	}

	return result, nil
}

// ReadHandler handles read requests for welcome
func ReadHandler(c *fiber.Ctx) error {
	// get welcome data and first page of thanks
	wdata, err := Get()
	thanksList, isLastPage, err := thanks.Read(1, 8)
	if err != nil {
		return server.APIInternalServerError(c)
	}

	// basic response data structure
	response := map[string]interface{}{
		"donated": map[string]interface{}{
			"count":  wdata.Count,
			"amount": wdata.Amount,
		},
		"thanks": map[string]interface{}{
			"results": thanksList,
		},
	}

	// changes depending to the last page
	if isLastPage {
		response["thanks"].(map[string]interface{})["_last"] = true
	} else {
		response["thanks"].(map[string]interface{})["_last"] = false
		response["thanks"].(map[string]interface{})["_next"] = "/thanks/2/?perPage=8"
	}

	return server.APIOK(c, "Požadavek byl úspěšně zpracován.", response)
}
