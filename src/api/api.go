package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/handlers"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/pkgs/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ApiServer struct {
	addr string
}

func NewApiServer(addr string) *ApiServer {

	return &ApiServer{
		addr: addr,
	}
}

func (s ApiServer) Start() {
	r := chi.NewRouter()
	setupMiddlewares(r)
	db := db.ConnectToDB()
	repositories.SetUpRepositories(db)

	r.Route("/api/v1", func(route chi.Router) {
		handlers := handlers.NewHandlers(&route)
		handlers.SetupRoutes(db)
	})

	fmt.Printf("Starting server...%s\n", s.addr)
	error := http.ListenAndServe(s.addr, r)

	if error != nil {
		log.Fatalf("Error starting server: %v", error)
		panic(error)
	}

}

func setupMiddlewares(r *chi.Mux) {
	router := *r
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
}
