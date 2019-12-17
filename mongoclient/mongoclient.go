package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DB_URI = "mongodb://localhost:27017"

// Create DB Connection
var ctx, client = initDB()

func initDB() (ctx context.Context, client *mongo.Client) {
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	return ctx, client
}
