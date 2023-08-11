package log

import (
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/pkg/configs/logging"
)

var (
	// Logger is a shared go-kit logger.
	// TODO: Change all components to take a non-global logger via their constructors.
	// Prefer accepting a non-global logger as an argument.
	Logger = log.NewNopLogger()
)

// InitLogger initialises the global gokit logger (util_log.Logger) and overrides the
// default logger for the server.
func InitLogger(cfg *logging.Config) {
	l := newBasicLogger(cfg.Format)

	// when using util_log.Logger, skip 5 stack frames.
	logger := log.With(l, "caller", log.Caller(5))
	// Must put the level filter last for efficiency.
	Logger = level.NewFilter(logger, level.Allow(cfg.LogLevel))
}

func newBasicLogger(format string) log.Logger {
	var logger log.Logger
	if format == "json" {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	} else {
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	}

	// return a Logger without filter or caller information, shouldn't use directly
	return log.With(logger, "ts", log.DefaultTimestampUTC)
}

// NewDefaultLogger creates a new gokit logger with the configured level and format
func NewDefaultLogger(l level.Value, format string) log.Logger {
	logger := newBasicLogger(format)
	return level.NewFilter(log.With(logger, "ts", log.DefaultTimestampUTC), level.Allow(l))
}

// CheckFatal prints an error and exits with error code 1 if err is non-nil
func CheckFatal(location string, err error) {
	if err != nil {
		logger := level.Error(Logger)
		if location != "" {
			logger = log.With(logger, "msg", "error "+location)
		}
		// %+v gets the stack trace from errors using github.com/pkg/errors
		logger.Log("err", fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
