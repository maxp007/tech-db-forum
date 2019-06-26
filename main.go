package main

import (
	"fmt"
	"github.com/maxp007/tech-db-forum/database"
	"github.com/maxp007/tech-db-forum/router"
	"log"
	"net/http"
)

func main() {

	err := http.ListenAndServe(":5000", router.GetRouter()) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	err = database.Create()
	if err != nil {

		fmt.Print("database.Connect", err)
	}
	defer func() {

		err := database.SQLConnection.Close()
		if err != nil {

			fmt.Print("defer database.SQLConnection.Close", err)
		}
	}()
}
