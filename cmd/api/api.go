package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/YhomiAce/rest-api-pg/services/users"
	"github.com/YhomiAce/rest-api-pg/config"
	"github.com/go-chi/cors"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
	config *config.Config
}

func NewAPIServer(addr string, db *sql.DB, config *config.Config) *APIServer {
	return &APIServer{addr: addr, db: db, config: config}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()

	// cors
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	corsHandler := c.Handler(router)

// Register routes
	userStore := users.NewUserStore(s.db, s.config)
	userHandler := users.NewHandler(userStore, s.config)
	userHandler.RegisterRoutes(router.PathPrefix("/api").Subrouter())

	log.Printf("Starting API server on %s\n", s.addr)

	return http.ListenAndServe(s.addr, corsHandler)
}