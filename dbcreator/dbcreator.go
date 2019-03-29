package dbcreator

import (
	"github.com/maxp007/tech-db-1/connector"
	"fmt"
	"io/ioutil"
)

func Create()(err error) {
	db, err:=connector.Connect()
	if err!=nil{
		fmt.Println("connector module:create()",err)
		return
	}
	file, err:=ioutil.ReadFile("tech-db-1")

	_, err =db.Exec(string(file))

	if err!=nil{
		fmt.Println("connector module:create()",err)
		return
	}
return
}
