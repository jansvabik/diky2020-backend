package welcome

import (
	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/internal/thanks"
	"github.com/noltio/diky2020-backend/pkg/server"
)

// ReadHandler handles read requests for welcome
func ReadHandler(c *fiber.Ctx) error {
	// get first page of thanks
	data, isLastPage, err := thanks.Read(1, 8)
	if err != nil {
		return server.APIInternalServerError(c)
	}

	// basic response data structure
	response := map[string]interface{}{
		"donated": map[string]interface{}{
			"count":  15497,
			"amount": 1583300,
		},
		"thanks": map[string]interface{}{
			"results": data,
		},
	}

	// changes depending to the last page
	if isLastPage {
		response["thanks"].(map[string]interface{})["_last"] = true
	} else {
		response["thanks"].(map[string]interface{})["_last"] = false
		response["thanks"].(map[string]interface{})["_next"] = "/thanks/2/"
	}

	return server.APIOK(c, "Požadavek byl úspěšně zpracován.", response)
}
