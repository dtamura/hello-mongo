package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
)

// Message はメッセージに対応する構造体
type Message struct {
	APIVersion string    `json:"apiVersion" bson:"apiVersion"`
	Body       string    `json:"body" binding:"required"`
	Traceid    string    `bson:"traceId"`
	Timestamp  time.Time `bson:"timestamp"`
	CreatedBy  string    `bson:"createdBy"`
	UpdatedBy  string    `bson:"updatedBy"`
}

func postMessage(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c, "POST /messages")
	defer span.Finish()

	hostname, _ := os.Hostname()

	span.SetTag("hostname", hostname) // Tagに"hello-to"をセット

	sp := opentracing.StartSpan(
		"process message",
		opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	var m Message
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprintf("%v", err)})
		return
	}
	processMessage(c, &m)

	sp1 := opentracing.StartSpan(
		"insert message",
		opentracing.ChildOf(span.Context()))
	defer sp1.Finish()

	res, err := insertOneMessage(c, mongoDb, m)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save message"})
	}
	c.JSON(200, gin.H{"body": m.Body, "apiVersion": m.APIVersion, "id": res.InsertedID})

}

// 格納用にメッセージを加工
func processMessage(ctx context.Context, m *Message) {
	m.Timestamp = time.Now()
}

func insertOneMessage(ctx context.Context, db *mongo.Database, m Message) (*mongo.InsertOneResult, error) {
	collName := "messages"
	coll := db.Collection(collName)

	return coll.InsertOne(ctx, m)
}
