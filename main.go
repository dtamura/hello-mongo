package main

import (
	"context"
	"os"
	"time"

	"github.com/dtamura/hello-mongo/lib/log"
	"github.com/dtamura/hello-mongo/lib/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	opentracing "github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var mongoClient *mongo.Client
var mongoDb *mongo.Database

// Tracing
var tracer opentracing.Tracer
var logger log.Factory

func main() {
	// loggerの初期化
	logger1, _ := zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel))
	zapLogger := logger1.With(zap.String("service", "hello-mongo"))
	logger = log.NewFactory(zapLogger)

	// OpenTracingの初期化
	tracer, closer := tracing.Init("hello-mongo", logger)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer) // Jaeger tracer のグローバル変数を初期化

	// MongoDBの初期化
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

	// Ginの初期化
	r := gin.Default()

	// Middleware
	r.Use(ginhttp.Middleware(tracer)) // Tracing

	// Router
	r.GET("/healthz", healthz)
	r.POST("/messages", postMessage)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
