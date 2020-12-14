package thanks

import (
	"time"

	"github.com/noltio/diky2020-backend/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Donation is a specific donation data structure
type Donation struct {
	Target string `json:"target" bson:"target"`
	Amount uint   `json:"amount" bson:"amount"`
}

// Thanks is a thanks data structure
type Thanks struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Text     string             `json:"text" bson:"text"`
	Time     *time.Time         `json:"time" bson:"time"`
	Donation *Donation          `json:"donation" bson:"donation"`
}

// collection returns the collection for the current data
func collection() *mongo.Collection {
	return app.State.MongoClient.Database(app.State.Cfg.Database.Name).Collection("thanks")
}
