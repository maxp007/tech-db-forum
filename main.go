package main

import (
	"fmt"
	"github.com/maxp007/tech-db-forum/db_creator"
)

func main() {
	err:= db_creator.Create()
	if err!=nil{
		fmt.Println("Main:db_creator.Create()", err)
		return
	}
}
