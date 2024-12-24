package logs

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Now()
		log.Printf("Method: %s | Path: %s | Start Time: %s | Duration: %s\n",
			r.Method, r.URL.Path, startTime.Format(time.RFC3339), endTime.Sub(startTime))
	})
}
