package main

import (
	"fmt"

	"github.com/felix-Asante/pennyPilot-go-api/src/api"
	"github.com/felix-Asante/pennyPilot-go-api/src/configs"
)

func main() {
	server := api.NewApiServer(fmt.Sprintf(":%s", configs.GetEnv("PORT")))
	server.Start()
}
