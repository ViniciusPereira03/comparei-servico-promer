package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"main/config"
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

func PubNewProduct(mercadoProdutoId int64, userID string) error {
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
	_, err = rdb.Publish(ctx, "new_product", string(payload)).Result()
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem no Redis: %v", err)
	}

	return nil
}

func PubUpdateProduct(mercadoProdutoId int64, userID string) error {
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
	_, err = rdb.Publish(ctx, "update_product", string(payload)).Result()
	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem no Redis: %v", err)
	}

	return nil
}
