package cli

import (
	"os"

	isatty "github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new logger.
func NewLogger(debug bool) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	if !debug {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	isTerm := isatty.IsTerminal(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"
	if isTerm {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.Encoding = "console"
	} else {
		config.Encoding = "json"
	}
	return config.Build()
}
