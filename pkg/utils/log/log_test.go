package log_test

import (
	"flag"
	"os"
	"testing"
	"time"

	gokitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/pkg/configs/logging"
	"github.com/udmire/observability-operator/pkg/utils/log"
)

// Check that debug lines are correctly filtered out.
func ExampleInitLogger() {
	// Kludge a couple of things so we can do tests repeatably.
	saveStderr := os.Stderr
	os.Stderr = os.Stdout
	saveTimestamp := gokitlog.DefaultTimestampUTC
	gokitlog.DefaultTimestampUTC = gokitlog.TimestampFormat(
		func() time.Time { return time.Unix(0, 0).UTC() },
		time.RFC3339Nano,
	)

	cfg := logging.Config{Level: "info"}

	cfg.RegisterFlags(&flag.FlagSet{})
	log.InitLogger(&cfg)
	level.Info(log.Logger).Log("test", "1")
	level.Debug(log.Logger).Log("test", "2 - should not print")
	// Output:
	// ts=1970-01-01T00:00:00Z caller=log_test.go:30 level=info test=1

	os.Stderr = saveStderr
	gokitlog.DefaultTimestampUTC = saveTimestamp
}

// Check the overhead of debug logging which gets filtered out.
func BenchmarkDebugLog(b *testing.B) {
	cfg := logging.Config{Level: "info"}
	cfg.RegisterFlags(&flag.FlagSet{})
	log.InitLogger(&cfg)
	b.ResetTimer()
	dl := level.Debug(log.Logger)
	for i := 0; i < b.N; i++ {
		dl.Log("something", "happened")
	}
}
