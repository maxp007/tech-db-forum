package main

import (
	"fmt"
	"github.com/maxp007/tech-db-forum/database"
	"github.com/maxp007/tech-db-forum/router"
	"log"
	"net/http"
)

func main() {
	log.Println("Started listening on port", 5000)
	fmt.Println("Started listening on port", 5000)
	err := http.ListenAndServe(":5000", router.GetRouter()) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		log.Println("Started listening on port", 5000)
	}
	/*
		err = database.Create()
		if err != nil {

			fmt.Print("database.Connect", err)
		} else {
			log.Println("Created database schema")
		}
	*/
	defer func() {
		err := database.SQLConnection.Close()
		if err != nil {
			fmt.Print("defer database.SQLConnection.Close", err)
		}
	}()
}
