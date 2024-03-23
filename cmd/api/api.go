package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/DracoR22/Red-Raccoon/internal/cart"
	"github.com/DracoR22/Red-Raccoon/internal/order"
	"github.com/DracoR22/Red-Raccoon/internal/product"
	"github.com/DracoR22/Red-Raccoon/internal/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// User route
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// Product route
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	// Order route
	orderStore := order.NewStore(s.db)

	// Cart route
	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
