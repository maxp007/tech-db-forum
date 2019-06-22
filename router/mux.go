package router

import (
	"github.com/gorilla/mux"
	"github.com/maxp007/tech-db-forum/handlers"
)

func GetRouter() (router *mux.Router) {

	router = mux.NewRouter()

	forum := router.PathPrefix("/forum").Subrouter()
	forum.HandleFunc("/create", handlers.PostForumCreate).Methods("POST")
	forum.HandleFunc("/{slug}/create", handlers.PostThreadCreate).Methods("POST")
	forum.HandleFunc("/{slug}/details", handlers.GetForumDetails).Methods("GET")
	forum.HandleFunc("/{slug}/threads", handlers.GetForumThreads).Methods("GET")
	forum.HandleFunc("/{slug}/users", handlers.GetForumUsers).Methods("GET")

	post := router.PathPrefix("/post").Subrouter()
	post.HandleFunc("/{id}/details", handlers.GetPostDetails).Methods("GET")
	post.HandleFunc("/{id}/details", handlers.PostPostUpdate).Methods("POST")

	service := router.PathPrefix("/service").Subrouter()
	service.HandleFunc("/clear", handlers.PostServiceClear).Methods("POST")
	service.HandleFunc("/status", handlers.GetServiceStatus).Methods("GET")

	thread := router.PathPrefix("/thread").Subrouter()
	thread.HandleFunc("/{slug_or_id}/create", handlers.PostPostCreate).Methods("POST")
	thread.HandleFunc("/{slug_or_id}/details", handlers.GetThreadInfo).Methods("GET")
	thread.HandleFunc("/{slug_or_id}/details", handlers.PostThreadUpdate).Methods("POST")
	thread.HandleFunc("/{slug_or_id}/posts", handlers.GetThreadPosts).Methods("GET")
	thread.HandleFunc("/{slug_or_id}/vote", handlers.PostThreadVote).Methods("POST")

	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{nickname}/create", handlers.PostUserCreate).Methods("POST")
	user.HandleFunc("/{nickname}/profile", handlers.GetUserProfile).Methods("GET")
	user.HandleFunc("/{nickname}/profile", handlers.PostUserUpdate).Methods("POST")

	return
}
