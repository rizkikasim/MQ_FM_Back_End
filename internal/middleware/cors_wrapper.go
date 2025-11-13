package middleware

import "net/http"

type corsMux struct {
	mux *http.ServeMux
}

func (c *corsMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	CORSMiddleware(c.mux).ServeHTTP(w, r)
}

func WrapWithCORS(mux *http.ServeMux) http.Handler {
	return &corsMux{mux: mux}
}
