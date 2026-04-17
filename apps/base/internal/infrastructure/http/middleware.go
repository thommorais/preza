package http

import (
	"crypto/rand"
	"encoding/hex"
	"preza/internal/domain/logger"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

func requestMiddleware(log logger.Logger) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		requestID := newRequestID()
		start := time.Now()

		e.Request.Header.Set("X-Request-Id", requestID)

		err := e.Next()

		status := 200
		if err != nil {
			status = 500
		}

		fields := []logger.Field{
			logger.F("request_id", requestID),
			logger.F("method", e.Request.Method),
			logger.F("path", e.Request.URL.Path),
			logger.F("status", status),
			logger.F("latency_ms", time.Since(start).Milliseconds()),
		}

		if err != nil {
			log.Error("request failed", fields...)
		} else {
			log.Info("request completed", fields...)
		}

		return err
	}
}

func newRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
