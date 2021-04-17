package internal

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

// NewLogger returns a new default logger that emits logs in a format
// compatible with SPIRE.
func NewLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
}
