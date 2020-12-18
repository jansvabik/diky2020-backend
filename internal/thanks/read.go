package thanks

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Read reads the specified amount of thanks from database
// and returns them in a slice of Thank structures
func Read(page int, perPage int, sortBy string, sortType int) ([]Thanks, bool, error) {
	// create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// mongodb options
	skip := int64((page - 1) * perPage)
	limit := int64(perPage + 1)
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort: map[string]int{
			sortBy: sortType,
		},
	}

	// get data from database
	cursor, err := collection().Find(ctx, bson.M{}, &opts)
	if err != nil {
		return nil, false, err
	}
	defer cursor.Close(ctx)

	// create an array of recipes and return them
	thanksList := []Thanks{}
	for cursor.Next(ctx) {
		var thanks Thanks
		cursor.Decode(&thanks)
		thanksList = append(thanksList, thanks)
	}

	// is this page the latest
	isLastPage := false
	var retList []Thanks
	if int64(len(thanksList)) != limit {
		isLastPage = true
		retList = thanksList
	} else {
		retList = thanksList[:len(thanksList)-1]
	}

	return retList, isLastPage, nil
}

// ReadHandler handles read requests for thanks
func ReadHandler(c *fiber.Ctx) error {
	// extract and validate the page param
	pageParam := c.Query("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		return server.APIError(c, "Číslo stránky musí být kladné celé číslo", 400)
	}

	// extract and validate the perPage param
	perPageParam := c.Query("perPage", "8")
	perPage, err := strconv.Atoi(perPageParam)
	if err != nil || perPage < 1 {
		return server.APIError(c, "Počet záznamů na stránce musí být kladné celé číslo", 400)
	}

	// extract and validate the sortBy param
	sortBy := c.Query("sortBy", "time")
	if sortBy != "time" && sortBy != "likes" {
		return server.APIError(c, "Řadit je možné jen podle času (time) a počtu liků (likes).", 400)
	}

	// extract and validate the sortType param
	sortTypeParam := c.Query("sortType", "desc")
	sortType := 1
	if sortTypeParam == "desc" {
		sortType = -1
	}

	// get thanks from database
	list, isLastPage, err := Read(page, perPage, sortBy, sortType)
	if err != nil {
		return server.APIInternalServerError(c)
	}

	// create response with the data and helper fields
	result := map[string]interface{}{
		"results": list,
		"_last":   false,
	}

	// changes for the latest page
	if isLastPage {
		result["_last"] = true
	} else {
		result["_next"] = "/thanks/" +
			"?page=" + strconv.Itoa(page+1) +
			"&perPage=" + strconv.Itoa(perPage) +
			"&sortBy=" + sortBy +
			"&sortType=" + sortTypeParam
	}

	return server.APIOK(c, "Požadavek byl úspěšně zpracován.", result)
}
