package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoDb *mongo.Database

func main() {
	r := gin.Default()
	mongoURL := os.Getenv("MONGO_URL")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	var err error
	mongoClient, err = mongo.NewClient(options.Client().ApplyURI(mongoURL))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Connect(ctx)
	if err != nil {
		return
	}
	mongoDb = mongoClient.Database(mongoDatabase)
	r.GET("/healthz", healthz)
	r.POST("/message", postMessage)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
