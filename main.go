package main

import (
	"fmt"
	"github.com/maxp007/tech-db-forum/router"
	"net/http"
)

func main() {
	fmt.Println("Started listening on port", 5000)
	fmt.Println("Started listening on port", 5000)
	err := http.ListenAndServe(":5000", router.GetRouter()) // задаем слушать порт

	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	} else {
		fmt.Println("Started listening on port", 5000)
	}
}
