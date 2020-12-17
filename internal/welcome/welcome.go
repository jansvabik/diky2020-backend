package welcome

import (
	"time"

	"github.com/noltio/diky2020-backend/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// WData is a data structure for homepage analytics data
type WData struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Amount   int                `json:"amount" bson:"amount"`
	Count    int                `json:"count" bson:"count"`
	EventEnd *time.Time         `json:"eventEnd" bson:"eventEnd"`
}

// collection returns the collection for the current data
func collection() *mongo.Collection {
	return app.State.MongoClient.Database(app.State.Cfg.Database.Name).Collection("welcome")
}
