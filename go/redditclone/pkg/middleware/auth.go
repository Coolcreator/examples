package middleware

import (
	"context"
	"log"

	"myapp/pkg/session"

	"net/http"

	"github.com/gorilla/mux"
)

var (
	noAuthUrls = map[string]bool{
		"/":             true,
		"/api/register": true,
		"/api/login":    true,
		"/api/posts/":   true,
		"/static/":      true,
	}
)

// func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Authentication middleware on path", r.URL.Path)
// 		if _, ok := noAuthUrls[r.URL.Path]; ok {
// 			log.Println("The path", r.URL.Path, "doesn't need auth")
// 			next.ServeHTTP(w, r)
// 			return
// 		}
// 		sess, err := sm.Check(r)
// 		log.Println("Session = ", sess, "error = ", err)
// 		if err != nil {
// 			http.Redirect(w, r, "/", 302)
// 			return
// 		}
// 		ctx := context.WithValue(r.Context(), session.SessionKey, sess)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func AuthWrapper(sm *session.SessionsManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Authentication middleware on path", r.URL.Path)
			path, err := mux.CurrentRoute(r).GetPathTemplate()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			log.Println(path)
			if ok, _ := noAuthUrls[path]; ok {
				log.Println("The path", r.URL.Path, "doesn't need auth")
				next.ServeHTTP(w, r)
				return
			}
			sess, err := sm.Check(r)
			log.Println("Session = ", sess, "error = ", err)
			if err != nil {
				http.Redirect(w, r, "/", 302)
				return
			}
			ctx := context.WithValue(r.Context(), session.SessionKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
