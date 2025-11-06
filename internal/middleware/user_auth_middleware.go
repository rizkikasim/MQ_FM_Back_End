package middleware

import (
	"errors"
	"net/http"
	"strings"

	authuserservice "mqfm_backend/internal/application/services/auth/user"
	errorinterceptor "mqfm_backend/internal/presentation/interceptor/errors"

)

// ---------------- USER AUTH MIDDLEWARE ----------------
// memastikan hanya user terautentikasi yang bisa akses endpoint.
func UserAuthMiddleware(authSvc *authuserservice.AuthUserService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorinterceptor.ErrorInterceptor(w, r, errors.New("token tidak ditemukan"), http.StatusUnauthorized, 0)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			errorinterceptor.ErrorInterceptor(w, r, errors.New("format token tidak valid"), http.StatusUnauthorized, 0)
			return
		}

		_, err := authSvc.Me(token)
		if err != nil {
			errorinterceptor.ErrorInterceptor(w, r, errors.New("token tidak valid atau kadaluarsa"), http.StatusUnauthorized, 0)
			return
		}

		// ✅ token valid → lanjut ke handler berikutnya
		next(w, r)
	}
}
