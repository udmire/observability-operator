package log

import (
	"github.com/go-kit/log"
)

// WithTraceID returns a Logger that has information about the traceID in
// its details.
func WithTraceID(traceID string, l log.Logger) log.Logger {
	// See note in WithContext.
	return log.With(l, "traceID", traceID)
}

// WithSourceIPs returns a Logger that has information about the source IPs in
// its details.
func WithSourceIPs(sourceIPs string, l log.Logger) log.Logger {
	return log.With(l, "sourceIPs", sourceIPs)
}
