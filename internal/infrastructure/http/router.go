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

	// 🔓 ROTAS PÚBLICAS (sem API key)
	r.HandleFunc("/produto/image/{barcode}", GetImageByBarcode).Methods("GET")

	// 🔒 ROTAS PROTEGIDAS
	api := r.PathPrefix("/").Subrouter()
	api.Use(middleware.APIKeyMiddleware)

	api.HandleFunc("/produto", CreateProduct).Methods("POST")
	api.HandleFunc("/produto", UpdateProduct).Methods("PUT")
	api.HandleFunc("/produto/{product_id}/{market_id}", GetProductByMarket).Methods("GET")
	api.HandleFunc("/produto/identificar", IdentifyProduct).Methods("POST")
	api.HandleFunc("/mercado", CreateMarket).Methods("POST")
	api.HandleFunc("/mercados", GetMarketByCoordinates).Methods("GET")
	api.HandleFunc("/mercado/produto/confirmar", ConfirmarValor).Methods("POST")
	api.HandleFunc("/produto/barcode/{barcode}", SearchProductByBarCode).Methods("GET")
	api.HandleFunc("/produto/search", SearchProductsByText).Methods("GET")

	return r
}
