package pocketbase_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func newTestApp(t *testing.T) *tests.TestApp {
	t.Helper()

	_, currentFile, _, _ := runtime.Caller(0)
	// currentFile: .../apps/base/internal/infrastructure/pocketbase/testhelper_test.go
	// pb_data:     .../apps/base/pb_data
	pbDataDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "pb_data")

	app, err := tests.NewTestApp(pbDataDir)
	if err != nil {
		t.Fatalf("create test app: %v", err)
	}

	t.Cleanup(app.Cleanup)
	return app
}

func bg() context.Context {
	return context.Background()
}
