package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"main/config"
	"main/internal/app"
	"main/internal/domain/user"
	"os"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var user_service *app.UserService

func init() {
	config.LoadConfig()
	host := fmt.Sprintf("%v:%v", os.Getenv("REDIS_MESSAGING_HOST"), os.Getenv("REDIS_MESSAGING_PORT"))
	rdb = redis.NewClient(&redis.Options{
		Addr: host,
	})
}

func SetUserService(service *app.UserService) {
	user_service = service
}

func Run() {
	go subCreateUser()
}

func subCreateUser() error {
	ctx := context.Background()

	sub := rdb.Subscribe(ctx, "user_created")
	ch := sub.Channel()

	for msg := range ch {
		var user user.User
		err := json.Unmarshal([]byte(msg.Payload), &user)
		if err != nil {
			fmt.Println("[ERRO] Erro ao decodificar payload de mensageria:", err)
			continue
		}

		err_create := user_service.CreateUser(&user)
		if err_create != nil {
			fmt.Println("[ERRO] Erro ao criar user:", err_create)
		}
	}

	return nil
}
