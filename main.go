package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client
var mongoDb *mongo.Database

func main() {
	r := gin.Default()
	mongoURL := os.Getenv("MONGO_URL")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	mongoDb = client.Database(mongoDatabase)
	mongoClient = client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}
	r.GET("/healthz", healthz)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
