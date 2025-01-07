package main

import (
	"fmt"

	"github.com/felix-Asante/pennyPilot-go-api/src/api"
	"github.com/felix-Asante/pennyPilot-go-api/src/utils"
)

func main() {
	server := api.NewApiServer(fmt.Sprintf(":%s", utils.GetEnv("PORT")))
	server.Start()
}
