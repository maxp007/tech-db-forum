package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "945020"
	dbname   = "tech-db-1"
)

func Connect(query string) (result sql.Result, err error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return
	}

	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}()

	err = db.Ping()
	if err != nil {
		return
	}

}

func main() {
	result, err := Connect(`SELECT * FROM  User`)
	if err != nil {
		fmt.Println("Execution error:", err)
		return
	}
	fmt.Println(result)

	/*PORT :=os.Getenv("PORT")
	if PORT == ""{
		PORT="5000"
	}
	r := router.GetRouter()
	err:=http.ListenAndServe(":"+PORT, r)
	if err!=nil{}
	fmt.Println("",err)*/

}
