package main

import (
	"fmt"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/db"
)

func main() {
	db := db.ConnectToDB()
	error := db.AutoMigrate(&repositories.Users{})

	if error != nil {
		panic(error)
	}
	fmt.Println("Migration successful")
}
