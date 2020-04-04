package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

func healthz(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c, "healthz")
	defer span.Finish()

	hostname, _ := os.Hostname()

	span.SetTag("hostname", hostname) // Tagに"hello-to"をセット

	sp := opentracing.StartSpan(
		"ping mongo",
		opentracing.ChildOf(span.Context()))
	defer sp.Finish()

	err := mongoClient.Connect(c)
	if err != nil && err != topology.ErrTopologyConnected {
		c.JSON(500, gin.H{
			"timestamp": time.Now(),
			"status":    "NG",
			"message":   "i'm not healthy",
			"error":     err,
			"hostname":  hostname,
		})
	}
	err = mongoClient.Ping(c, readpref.Primary())
	if err != nil {
		c.JSON(500, gin.H{
			"timestamp": time.Now(),
			"status":    "NG",
			"message":   "i'm not healthy",
			"error":     err,
			"hostname":  hostname,
		})

	}
	c.JSON(200, gin.H{
		"timestamp": time.Now(),
		"status":    "OK",
		"message":   "i'm healthy",
		"hostname":  hostname,
	})
}
