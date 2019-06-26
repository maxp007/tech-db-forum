package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "max007"
	password = "12345qwerty"
	dbname   = "tech-db-1"
)

var SQLConnection *sql.DB

func Connect() (db *sql.DB, err error) {

	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", connString)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	if err != nil {
		fmt.Print("Connection open error", err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Print("Connectin ping", err)
	}
	return
}

func init() {
	conn, err := Connect()
	SQLConnection = conn

	if err != nil {
		fmt.Println(err)
	}
}
