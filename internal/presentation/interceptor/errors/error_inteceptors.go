package errors

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/google/uuid"

)

// ErrorResponse bentuk JSON untuk response API error
type ErrorResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	AdminID   int    `json:"admin_id"`
	TraceID   string `json:"trace_id"`
}

// ErrorInterceptor menulis log error di terminal & kirim response JSON ke client
func ErrorInterceptor(w http.ResponseWriter, r *http.Request, err error, status int, adminID int) {
	if err == nil {
		return
	}

	traceID := uuid.New().String() // tetap generate trace_id untuk cross-reference log

	// ambil detail lokasi error
	_, file, line, _ := runtime.Caller(1)

	// setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	// log ke terminal
	logger.Error("request failed",
		slog.Int("admin_id", adminID),
		slog.String("trace_id", traceID),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("error", err.Error()),
		slog.String("file", file),
		slog.Int("line", line),
	)

	// response ke client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Status:    status,
		Message:   err.Error(),
		Timestamp: time.Now().Format(time.RFC3339),
		AdminID:   adminID,
		TraceID:   traceID,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
