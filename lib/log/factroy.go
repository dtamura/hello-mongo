package log

import "go.uber.org/zap"

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type Factory struct {
	logger *zap.Logger
}

// NewFactory creates a new Factory.
func NewFactory(logger *zap.Logger) Factory {
	return Factory{logger: logger}
}

// Bg creates a context-unaware logger.
func (b Factory) Bg() Logger {
	return logger(b)
}
