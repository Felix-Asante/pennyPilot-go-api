package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/felix-Asante/pennyPilot-go-api/src/api/handlers"
	"github.com/felix-Asante/pennyPilot-go-api/src/api/repositories"
	"github.com/felix-Asante/pennyPilot-go-api/src/configs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db := createDbAndRepositories()
	r.Route("/api/v1", func(route chi.Router) {
		handlers := handlers.NewHandlers(&route)
		handlers.SetupRoutes(db)
	})

	fmt.Printf("Starting server...%s", s.addr)
	error := http.ListenAndServe(s.addr, r)

	if error != nil {
		log.Fatalf("Error starting server: %v", error)
		panic(error)
	}

}

func createDbAndRepositories() *gorm.DB {
	db_host := configs.GetEnv("DB_HOST")
	db_port := configs.GetEnv("DB_PORT")
	db_user := configs.GetEnv("DB_USER")
	db_name := configs.GetEnv("DB_NAME")

	dsn := fmt.Sprintf("host=%v user=%v password= dbname=%v port=%v sslmode=disable", db_host, db_user, db_name, db_port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database%v", err)
		panic("failed to connect database")
	}
	fmt.Println("Connected to database")
	repositories.SetUpRepositories(db)
	return db
}

func setupMiddlewares(r *chi.Mux) {
	router := *r
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
}
