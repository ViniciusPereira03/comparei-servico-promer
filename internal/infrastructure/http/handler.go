package http

import (
	"encoding/json"
	"fmt"
	"main/internal/app"
	"main/internal/infrastructure/http/dto"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

var service *app.ProdutosService
var mercado_service *app.MercadoService

func IniHandlers(produtosService *app.ProdutosService, mercadoService *app.MercadoService) {
	service = produtosService
	mercado_service = mercadoService
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, err error, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":    err.Error(),
		"mensagem": message,
	})
}

func validaToken(w http.ResponseWriter, r *http.Request) (string, error) {
	secret := os.Getenv("USER_JWT_SECRET")

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("Missing token")
	}

	// Remover o prefixo "Bearer " se existir
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Verificar o token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	// Acessar os dados (claims)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id := claims["id"]
		return fmt.Sprintf("%v", id), nil
	}

	return "", fmt.Errorf("Erro ao decodificar token")
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	userID, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	var produtoDTO dto.CreateProductDTO
	if err := json.NewDecoder(r.Body).Decode(&produtoDTO); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "JSON inválido")
		return
	}

	produto := produtoDTO.ParseToProduct()
	id, err := service.CreateProduct(produto, userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Erro ao refistrar log")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fmt.Sprintf("Produto %v cadastrado com sucesso!", id))
}

func IdentifyProduct(w http.ResponseWriter, r *http.Request) {
	_, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	var body dto.ProdutoFotoDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "JSON inválido")
		return
	}

	produto, err := service.IdetificarProduto(body.ParseToPhoto())
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err, "Erro ao refistrar log")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(produto)

}

func GetMarketByCoordinates(w http.ResponseWriter, r *http.Request) {
	userID, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	user, err_user := user_service.GetUser(userID)
	if err_user != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_user, "Erro ao refistrar log")
		return
	}

	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	fmt.Printf(latStr)
	fmt.Printf(lngStr)

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Latitude inválida", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "Longitude inválida", http.StatusBadRequest)
		return
	}

	mercados, err := mercado_service.GetMarketByCoordinates(lat, lng, user.RayDistance)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "JSON inválido")
	}

	json.NewEncoder(w).Encode(mercados)
}

func ConfirmarValor(w http.ResponseWriter, r *http.Request) {
	userID, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	var body dto.ConfirmaMercadoProdutoDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "JSON inválido")
		return
	}

	_, err := mercado_service.ConfirmarValor(body.ParseToMercadoProdutos(), userID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "JSON inválido")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Confirmado")
}

func SearchProductByBarCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	barcode := vars["barcode"]

	mercados, err := service.SearchProductsByBarcode(barcode)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, err.Error())
		return
	}

	json.NewEncoder(w).Encode(mercados)
}

func SearchProductsByText(w http.ResponseWriter, r *http.Request) {
	userID, err_token := validaToken(w, r)
	if err_token != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_token, "Erro ao refistrar log")
		return
	}

	user, err_user := user_service.GetUser(userID)
	if err_user != nil {
		sendErrorResponse(w, http.StatusInternalServerError, err_user, "Erro ao refistrar log")
		return
	}

	text := r.URL.Query().Get("text")
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Latitude inválida", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "Longitude inválida", http.StatusBadRequest)
		return
	}

	mercadosRegiao, err := mercado_service.SearchMarketsByCoordinates(lat, lng, user.RayDistance)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err, "JSON inválido")
	}

	if len(mercadosRegiao) > 0 {

		mercados, err := service.SearchProductsByText(text)
		if err != nil {
			sendErrorResponse(w, http.StatusBadRequest, err, err.Error())
			return
		}

		var mercadosFiltrados []*mercadoprodutos.MercadoProdutosCompleto
		for _, m := range mercados {
			for _, mr := range mercadosRegiao {
				if m.Mercado.ID == mr.ID {
					mercadosFiltrados = append(mercadosFiltrados, m)
					break
				}
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mercadosFiltrados)
	} else {
		http.Error(w, fmt.Sprintf("Nenhum mercado cadastrado no raio de %d Km", user.RayDistance), http.StatusBadRequest)
	}
}
}
