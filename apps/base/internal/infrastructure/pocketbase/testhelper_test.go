package pocketbase_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"preza/internal/domain/logger"

	"github.com/pocketbase/pocketbase/tests"
)

func newTestApp(t *testing.T) *tests.TestApp {
	t.Helper()

	_, currentFile, _, _ := runtime.Caller(0)
	pbDataDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "pb_data")

	app, err := tests.NewTestApp(pbDataDir)
	if err != nil {
		t.Fatalf("create test app: %v", err)
	}

	t.Cleanup(app.Cleanup)
	return app
}

// noopLogger satisfies logger.Logger in tests without any output.
type noopLogger struct{}

func (noopLogger) Info(_ string, _ ...logger.Field)  {}
func (noopLogger) Warn(_ string, _ ...logger.Field)  {}
func (noopLogger) Error(_ string, _ ...logger.Field) {}

func testLogger() logger.Logger { return noopLogger{} }

func bg() context.Context {
	return context.Background()
}
