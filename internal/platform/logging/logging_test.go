package logging_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"go-market/internal/platform/logging"
)

func TestNewAddsServiceToLogRecords(t *testing.T) {
	var output bytes.Buffer
	logger := logging.New(&output, "catalog", slog.LevelInfo)

	logger.Info("service started")

	if !strings.Contains(output.String(), `"service":"catalog"`) {
		t.Errorf("log record does not contain service: %s", output.String())
	}
}
