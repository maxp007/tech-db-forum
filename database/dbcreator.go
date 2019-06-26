package database

import (
	"fmt"
	"io/ioutil"
)

func Create() (err error) {
	log.Println("schema Creator method")
	db, err := Connect()
	if err != nil {
		fmt.Println("database module:create() Connect", err)
		return
	}
	file, err := ioutil.ReadFile("dump.sql")
	if err != nil {
		fmt.Println("ioutil.ReadFile dump.sql", err)
		return
	}
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
	log.Println("Exitting Creator")

	return
}
