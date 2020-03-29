package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func healthz(c *gin.Context) {
	hostname, _ := os.Hostname()

	err := mongoClient.Ping(c, readpref.Primary())
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
