package main

import (
	"bytes"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HTTPWriter is a custom writer that sends logs to an HTTP endpoint.
type HTTPWriter struct {
	URL string
}

// Write sends the log entry to the configured HTTP endpoint.
func (h *HTTPWriter) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", h.URL, bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, err
	}

	return len(p), nil
}

func main() {
	httpWriter := &HTTPWriter{
		URL: "http://localhost:3000/",
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // e.g. INFO, WARN
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // e.g. 2021-01-01T00:00:00Z
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(httpWriter),
		zapcore.InfoLevel,
	)

	logger := zap.New(core).With(
		zap.String("service", "golang-client"),
	)

	logger.Info("This is an informational message")
	logger.Error("This is an error message")
}
