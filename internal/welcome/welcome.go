package welcome

import (
	"github.com/noltio/diky2020-backend/pkg/app"
	"go.mongodb.org/mongo-driver/mongo"
)

// collection returns the collection for the current data
func collection() *mongo.Collection {
	return app.State.MongoClient.Database(app.State.Cfg.Database.Name).Collection("welcome")
}
