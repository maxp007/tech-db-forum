package main

import (
	"github.com/maxp007/tech-db-forum/router"
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":5000", router.GetRouter()) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
