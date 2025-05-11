package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepository implementa o repository usando MongoDB
type MongoRepository struct {
	collection *mongo.Collection
}

// NewMongoRepository cria um novo MongoRepository
func NewMongoRepository(client *mongo.Client, dbName, collectionName string) *MongoRepository {
	coll := client.Database(dbName).Collection(collectionName)
	return &MongoRepository{collection: coll}
}

func (r *MongoRepository) SaveImage(base64Image string, nome string) error {
	if base64Image == "" || nome == "" {
		return errors.New("imagem ou nome inválido")
	}

	// Separar o prefixo (tipo MIME) do conteúdo
	parts := strings.SplitN(base64Image, ",", 2)
	if len(parts) != 2 {
		return errors.New("formato base64 inválido")
	}

	mimePart := parts[0]
	dataPart := parts[1]

	// Verifica tipo de imagem
	var mimeType string
	if strings.Contains(mimePart, "image/png") {
		mimeType = "image/png"
	} else if strings.Contains(mimePart, "image/jpeg") {
		mimeType = "image/jpeg"
	} else {
		return errors.New("tipo de imagem não suportado")
	}

	// Decodifica o conteúdo base64
	imageBytes, err := base64.StdEncoding.DecodeString(dataPart)
	if err != nil {
		return err
	}

	// Cria um documento com os dados da imagem
	doc := bson.M{
		"_id":      nome, // Usar nome como ID evita duplicação
		"image":    imageBytes,
		"mimeType": mimeType,
		"created":  time.Now(),
	}

	// Substitui se já existir
	opts := options.Replace().SetUpsert(true)
	_, err = r.collection.ReplaceOne(context.Background(), bson.M{"_id": nome}, doc, opts)
	return err
}

func (r *MongoRepository) GetImageByName(nome string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": nome}

	var result struct {
		ID    string `bson:"_id"`
		Image string `bson:"image"`
	}

	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", errors.New("imagem não encontrada no banco de dados")
	}

	return result.Image, nil
}
