package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "945020"
	dbname   = "tech-db-1"
)

func Connect() (db *sql.DB, err error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", connString)
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
