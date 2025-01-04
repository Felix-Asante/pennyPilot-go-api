package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Felix-Asante/pennyPilot-go-api/src/api/handlers"
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
	setupMiddlewares(r)
	r.Route("/api/v1", func(route chi.Router) {
		handlers := handlers.NewHandlers(&route)
		handlers.SetupRoutes()
	})

	fmt.Printf("Starting server...%s", s.port)
	error := http.ListenAndServe(s.port, r)

	if error != nil {
		log.Fatalf("Error starting server: %v", error)
		panic(error)
	} else {
		fmt.Print("Server running on port " + s.port)
	}

}

func setupMiddlewares(r *chi.Mux) {
	router := *r
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
}
