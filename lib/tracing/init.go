package tracing

import (
	"fmt"
	"io"

	"github.com/dtamura/hello-mongo/lib/log"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

// Init returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func Init(serviceName string, logger log.Factory) (opentracing.Tracer, io.Closer) {
	cfg, err := config.FromEnv()
	if err != nil {
		logger.Bg().Fatal("cannot parse Jaeger env vars", zap.Error(err))
	}
	cfg.ServiceName = serviceName
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1
	cfg.Reporter.LogSpans = true

	jaegerLogger := jaegerLoggerAdapter{logger.Bg()}
	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerLogger),
	)
	if err != nil {
		logger.Bg().Fatal("ERROR: cannot init Jaeger: ", zap.Error(err))
	}
	return tracer, closer
}

type jaegerLoggerAdapter struct {
	logger log.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.logger.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, args...))
}
