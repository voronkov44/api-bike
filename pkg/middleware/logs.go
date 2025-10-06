package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)

		elapsed := time.Since(start)
		email := wrapper.Email()
		if email == "" {
			email = "Unauthorized"
		}

		log.Printf(
			"| email=%s | %3d %s %-30s | %7.2fms",
			email,
			wrapper.StatusCode,
			r.Method,
			r.URL.Path,
			float64(elapsed.Microseconds())/1000,
		)
	})
}
