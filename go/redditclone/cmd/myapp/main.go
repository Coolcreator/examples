package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"myapp/pkg/handlers"
	"myapp/pkg/items"
	"myapp/pkg/middleware"
	"myapp/pkg/session"
	"myapp/pkg/user"
)

func main() {
	template := template.Must(template.ParseFiles("../../template/index.html"))

	sm := session.NewSessionsMem()
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	userRepo := user.NewUserRepo()
	itemsRepo := items.NewRepo()

	userHandler := &handlers.UserHandler{
		Tmpl:     template,
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sm,
	}

	handlers := &handlers.ItemsHandler{
		Tmpl:      template,
		Logger:    logger,
		ItemsRepo: itemsRepo,
	}

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("../../template/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", userHandler.Index).Methods("GET")

	r.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	r.HandleFunc("/api/logout", userHandler.Logout).Methods("POST")

	r.HandleFunc("/api/posts/", handlers.List).Methods("GET")
	r.HandleFunc("/api/posts", handlers.AddPost).Methods("POST")
	r.HandleFunc("/api/post/{id}", handlers.GetPost).Methods("GET")
	r.HandleFunc("/api/post/{id}", handlers.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{id}", handlers.DeletePost).Methods("DELETE")
	r.HandleFunc("/api/post/{id}/upvote", handlers.Upvote).Methods("GET")
	r.HandleFunc("/api/post/{id}/downvote", handlers.Downvote).Methods("GET")
	r.HandleFunc("/api/post/{id}/unvote", handlers.Unvote).Methods("GET")
	r.HandleFunc("/api/user/{username}", handlers.UserList).Methods("GET")
	r.HandleFunc("/api/posts/{category}", handlers.CategoryList).Methods("GET")
	r.HandleFunc("/api/post/{id}/{comment}", handlers.DeleteComment).Methods("DELETE")
	r.Use(middleware.Panic)
	r.Use(middleware.AuthWrapper(sm))

	// mux := middleware.Auth(sm, r)
	// mux = middleware.AccessLog(logger, mux)
	// mux = middleware.Panic(mux)

	addr := ":8081"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, r)
}
