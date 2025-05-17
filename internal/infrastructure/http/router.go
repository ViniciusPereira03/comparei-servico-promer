package http

import (
	"main/internal/app"
	"main/internal/infrastructure/http/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func NewRouter(produtosService *app.ProdutosService) *mux.Router {
	_ = godotenv.Load()

	r := mux.NewRouter()

	r.Use(middleware.APIKeyMiddleware)

	r.HandleFunc("/produto", CreateProduct).Methods("POST")
	r.HandleFunc("/produto/identificar", IdentifyProduct).Methods("POST")
	r.HandleFunc("/mercados", GetMarketByCoordinates).Methods("GET")

	return r
}
