package successinterceptor

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

)

// SuccessResponse bentuk JSON standar untuk response sukses
type SuccessResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
	AdminID   int         `json:"admin_id"`
	TraceID   string      `json:"trace_id"`
}

// SuccessInterceptor menulis log sukses di terminal & kirim response JSON ke client
func SuccessInterceptor(w http.ResponseWriter, r *http.Request, status int, message string, data interface{}, adminID int) {
	traceID := uuid.New().String()

	// setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// log ke terminal
	logger.Info("request success",
		slog.Int("admin_id", adminID),
		slog.String("trace_id", traceID),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", status),
	)

	// response ke client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := SuccessResponse{
		Status:    status,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
		AdminID:   adminID,
		TraceID:   traceID,
	}

	_ = json.NewEncoder(w).Encode(resp)
}
