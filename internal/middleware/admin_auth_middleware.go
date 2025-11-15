package middleware

import (
	"errors"
	"net/http"
	"strings"

	authadminservice "mqfm_backend/internal/application/services/auth/admin"
	errorinterceptor "mqfm_backend/internal/presentation/interceptor/errors"

)

/*
===================================================
🔥 ADMIN CORS MIDDLEWARE — FIX TOTAL NGROK + BROWSER
===================================================
*/
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Authorization, X-Requested-With, Accept, ngrok-skip-browser-warning",
		)

		w.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, OPTIONS")

		w.Header().Set("Access-Control-Expose-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}


/*
===================================================
🔥 ADMIN AUTH MIDDLEWARE (Token Guard)
===================================================
*/
func AdminAuthMiddleware(authSvc *authadminservice.AuthAdminService, next http.HandlerFunc) http.HandlerFunc {
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

		next(w, r)
	}
}
