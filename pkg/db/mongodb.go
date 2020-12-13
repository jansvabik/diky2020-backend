package db

import (
	"context"
	"net/url"
	"time"

	"github.com/noltio/tiskarnik-ms-auth/pkg/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConnect creates the database connection
func MongoConnect() (*mongo.Client, error) {
	// try to connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	options := options.Client().ApplyURI(
		"mongodb://" +
			app.State.Cfg.Database.User +
			":" +
			url.QueryEscape(app.State.Cfg.Database.Password) +
			"@" +
			app.State.Cfg.Database.Host +
			":" +
			app.State.Cfg.Database.Port +
			"/?authSource=" +
			app.State.Cfg.Database.Name +
			"&connect=direct")

	return mongo.Connect(ctx, options)
}
