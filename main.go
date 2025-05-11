package main

import (
	"context"
	"database/sql"
	"log"
	"main/config"
	"main/internal/app"
	customHTTP "main/internal/infrastructure/http"
	"main/internal/infrastructure/repository"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"net/http"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Erro ao carregar configurações:", err)
	}

	// Testar conexão com Redis de mensageria
	redisMessaging := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_MESSAGING_HOST") + ":" + os.Getenv("REDIS_MESSAGING_PORT"),
	})
	ctx := context.Background()
	_, err := redisMessaging.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Não foi possível conectar ao Redis de mensageria:", err)
	}

	// Configuração da conexão com o MySQL usando variáveis de ambiente
	dsn := os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ")/" + os.Getenv("MYSQL_DB") + "?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao abrir conexão MySQL:", err)
	}

	// Verificar a conexão com o MySQL
	if err := db.Ping(); err != nil {
		log.Fatal("Não foi possível conectar ao MySQL:", err)
	}

	// --- MongoDB ---
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("Erro ao criar cliente MongoDB:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := mongoClient.Connect(ctx); err != nil {
		log.Fatal("Erro ao conectar no MongoDB:", err)
	}
	// opcional: ping para certificar
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal("Ping no MongoDB falhou:", err)
	}

	mysqlRepo := repository.NewMySQLRepository(db)
	if mysqlRepo == nil {
		log.Fatal("mysqlRepo está nil")
	}
	mongoRepo := repository.NewMongoRepository(
		mongoClient,
		os.Getenv("MONGO_DB_NAME"),
		os.Getenv("MONGO_COLLECTION"),
	)
	mercadoService := app.NewMercadoService(mysqlRepo)
	productService := app.NewProductsService(mysqlRepo, mercadoService, mongoRepo)

	// Iniciar o servidor HTTP
	customHTTP.IniHandlers(productService, mercadoService)
	router := customHTTP.NewRouter(productService)

	// Inicia o servidor HTTP
	log.Println("🚀 Servidor rodando na porta " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}
