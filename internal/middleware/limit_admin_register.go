package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"

	authadminrepository "mqfm_backend/internal/infrastructure/repository/auth/admin"

)

// LimitAdminRegisterMiddleware membatasi jumlah admin maksimal 3
func LimitAdminRegisterMiddleware(repo authadminrepository.AuthAdminRepository, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		admins := repo.GetAll()
		traceID := uuid.New().String()

		if len(admins) >= 3 {
			// setup logger
			logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelWarn,
			}))

			logger.Warn("❌ pendaftaran admin ditolak (limit tercapai)",
				slog.String("trace_id", traceID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("total_admin", len(admins)),
				slog.Time("time", time.Now()),
			)

			// response ke client
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)

			resp := map[string]interface{}{
				"status":      http.StatusForbidden,
				"message":     "pendaftaran admin baru ditolak, batas maksimal 3 admin sudah tercapai",
				"trace_id":    traceID,
				"total_admin": len(admins),
				"timestamp":   time.Now().Format(time.RFC3339),
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		// lanjut ke handler berikutnya
		next.ServeHTTP(w, r)
	}
}
