package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"main/config"
	mercadoprodutos "main/internal/domain/mercado_produtos"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
	config.LoadConfig()
	host := fmt.Sprintf("%v:%v", os.Getenv("REDIS_MESSAGING_HOST"), os.Getenv("REDIS_MESSAGING_PORT"))
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
}

func PubNewProduct(mercadoProduto *mercadoprodutos.MercadoProdutos, userID string) error {
	ctx := context.Background()

	type payload_log struct {
		Id               int64   `json:"id"`
		MercadoProdutoId int64   `json:"mercado_produto_id"`
		UserID           string  `json:"user_id"`
		PrecoUnitario    float32 `json:"preco_unitario"`
		NivelConfianca   int32   `json:"nivel_confianca"`
		MercadoId        int64   `json:"mercado_id"`
		ProdutoID        int64   `json:"produto_id"`
	}

	var p payload_log
	p.Id = mercadoProduto.ID
	p.MercadoProdutoId = mercadoProduto.ID
	p.UserID = userID
	p.PrecoUnitario = mercadoProduto.PrecoUnitario
	p.NivelConfianca = mercadoProduto.NivelConfianca
	p.MercadoId = mercadoProduto.MercadoID
	p.ProdutoID = mercadoProduto.ProdutoID

	payload, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao codificar payload: %v", err)
	}
	_, err = rdb.Publish(ctx, "new_product", string(payload)).Result()
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem no Redis: %v", err)
	}

	return nil
}

func PubUpdateProduct(mercadoProduto *mercadoprodutos.MercadoProdutos, userID string) error {
	ctx := context.Background()

	type payload_log struct {
		Id             int64                            `json:"id"`
		UserID         string                           `json:"user_id"`
		MercadoProduto *mercadoprodutos.MercadoProdutos `json:"mercado_produto"`
	}

	var p payload_log
	p.Id = mercadoProduto.ID
	p.UserID = userID
	p.MercadoProduto = mercadoProduto

	payload, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao codificar payload: %v", err)
	}
	_, err = rdb.Publish(ctx, "update_product", string(payload)).Result()
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem no Redis: %v", err)
	}

	return nil
}

func PubConfirmaValor(mercadoProdutoId int64, userID string) error {
	ctx := context.Background()

	type payload_log struct {
		Id     int64  `json:"id"`
		UserID string `json:"user_id"`
	}

	var p payload_log
	p.Id = mercadoProdutoId
	p.UserID = userID

	payload, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao codificar payload: %v", err)
	}
	_, err = rdb.Publish(ctx, "confirma_valor_mercado_produto", string(payload)).Result() //FALTA CRIAR O SUB NO SERVIÇO DE LOG
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem no Redis: %v", err)
	}

	return nil
}
