package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ApiServer struct {
	port string
}

func NewApiServer(port string) *ApiServer {

	return &ApiServer{
		port: port,
	}
}

func (s ApiServer) Start() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fmt.Printf("Starting server...%s", s.port)
	error := http.ListenAndServe(s.port, r)

	if error != nil {
		log.Fatalf("Error starting server: %v", error)
		panic(error)
	} else {
		fmt.Print("Server running on port " + s.port)
	}

}
