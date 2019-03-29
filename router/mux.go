package router

import (
	"github.com/gorilla/mux"

)

func GetRouter() (router *mux.Router) {

	router = mux.NewRouter()

	forum := router.PathPrefix("/forum").Subrouter()
	forum.HandleFunc("/create", handlers.Connect).Methods("POST")
	forum.HandleFunc("/{slug}/create", handlers.Connect).Methods("POST")
	forum.HandleFunc("/{slug}/details", handlers.Connect).Methods("GET")
	forum.HandleFunc("/{slug}/threads", handlers.Connect).Methods("GET")
	forum.HandleFunc("/{slug}/users", handlers.Connect).Methods("GET")

	post := router.PathPrefix("/post").Subrouter()
	post.HandleFunc("/{id}/details", handlers.Connect).Methods("GET")
	post.HandleFunc("/{id}/details", handlers.Connect).Methods("POST")

	service := router.PathPrefix("/service").Subrouter()
	service.HandleFunc("/clear", handlers.Connect).Methods("POST")
	service.HandleFunc("/status", handlers.Connect).Methods("GET")

	thread := router.PathPrefix("/thread").Subrouter()
	thread.HandleFunc("/{slug_or_id}/create", handlers.Connect).Methods("POST")
	thread.HandleFunc("/{slug_or_id}/details", handlers.Connect).Methods("GET")
	thread.HandleFunc("/{slug_or_id}/details", handlers.Connect).Methods("POST")
	thread.HandleFunc("/{slug_or_id}/posts", handlers.Connect).Methods("GET")
	thread.HandleFunc("/{slug_or_id}/vote", handlers.Connect).Methods("POST")

	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{nickname}/create", handlers.Connect).Methods("POST")
	user.HandleFunc("/{nickname}/profile", handlers.Connect).Methods("GET")
	user.HandleFunc("/{nickname}/profile", handlers.Connect).Methods("POST")

	return
}
