package middleware

import (
	"log"
	"net/http"
)

func Panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Panic middleware on path", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
