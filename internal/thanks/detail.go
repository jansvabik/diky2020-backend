package thanks

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noltio/diky2020-backend/pkg/server"
	"go.mongodb.org/mongo-driver/bson"
)

func getThanks(sid string) (Thanks, error) {
	// create a background context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get data from database
	var thx Thanks
	err := collection().FindOne(ctx, bson.M{
		"shortId": sid,
	}).Decode(&thx)
	return thx, err
}

// DetailHandler handles requests for one specified thanks record
func DetailHandler(c *fiber.Ctx) error {
	// test the id param
	sid := c.Params("id")
	if len(sid) == 0 {
		return server.APIError(c, "Seš idiot, protože nezadals ID.", 400)
	}

	// get the document and send server response
	thx, err := getThanks(sid)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return server.APIError(c, "Poděkování s tímto shortId neexistuje.", 404)
		}
		return server.APIInternalServerError(c)
	}

	return server.APIOK(c, "Požadavek byl úspěšně zpracován.", thx)
}
