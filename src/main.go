package main

import (
	"fmt"

	"github.com/Felix-Asante/pennyPilot-go-api/src/api"
	"github.com/Felix-Asante/pennyPilot-go-api/src/configs"
)

func main() {
	server := api.NewApiServer(fmt.Sprintf(":%s", configs.PORT))
	server.Start()
}
