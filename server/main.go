package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mylogger "github.com/dtamura/hello-mongo/lib/log"
	"github.com/dtamura/hello-mongo/lib/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Name 名前
	Name string
	// Version バージョン
	Version string
	// Revision リビジョン
	Revision string
)

// AppConfig 読み込む設定の型
type AppConfig struct {
	Mode   string
	Server struct {
		Address string
		Port    string
	}
	MongoDB struct {
		URL      string
		Database string
	}
	Ping struct {
		URL string
	}
}

// mongodb
var appConfig AppConfig
var mongoClient *mongo.Client
var mongoDb *mongo.Database

// tracing
var tracer opentracing.Tracer
var logger mylogger.Factory

// Ping
var pingURL string

// Start the REST API server using the configuration provided
func Start(conf AppConfig) error {

	fmt.Printf("AppConfig: %v\n", conf)
	appConfig = conf
	if conf.Mode != "" {
		gin.SetMode(conf.Mode)
	}

	setupApp()

	router := SetupRouter()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", conf.Server.Address, conf.Server.Port),
		Handler: router,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// IdleTimeout:    120 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	finilize(ctx)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")

	return nil
}

// setupApp アプリケーション起動時処理
// MongoDBに接続
// OpenTracingの初期化
func setupApp() {
	// loggerの初期化
	logger1, _ := zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel))
	zapLogger := logger1.With(zap.String("service", "hello-mongo"))
	logger = mylogger.NewFactory(zapLogger)

	// OpenTracingの初期化
	var closer io.Closer
	tracer, closer = tracing.Init("hello-mongo", logger)
	defer closer.Close()

	var err error
	mongoClient, err = mongo.NewClient(options.Client().ApplyURI(appConfig.MongoDB.URL))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Connect(ctx)
	if err != nil {
		return
	}
	mongoDb = mongoClient.Database(appConfig.MongoDB.Database)
}

// finilize アプリケーション終了時の処理
func finilize(ctx context.Context) {
	err := mongoClient.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
