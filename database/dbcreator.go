package database

import (
	"fmt"
	"io/ioutil"
)

func Create() (err error) {
	db, err := Connect()
	if err != nil {
		fmt.Println("database module:create() Connect", err)
		return
	}
	file, err := ioutil.ReadFile("tech-db-1.sql")

	_, err = db.Exec(string(file))

	if err != nil {
		fmt.Println("database module:create() exec", err)
		return
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Print("Connector close error", err)
		}
	}()
	return
}
