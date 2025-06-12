package mercadoprodutos

import (
	"main/internal/domain/mercados"
	"main/internal/domain/produtos"
	"time"
)

type MercadoProdutos struct {
	ID             int64      `json:"id"`
	ProdutoID      int64      `json:"id_produto"`
	MercadoID      int64      `json:"id_mercado"`
	PrecoUnitario  float32    `json:"preco_unitario"`
	NivelConfianca int32      `json:"nivel_confianca"`
	CreatedAt      time.Time  `json:"created_at"`
	ModifiedAt     time.Time  `json:"modified_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type MercadoProdutosCompleto struct {
	ID             int64             `json:"id_mercado_produtos"`
	Produto        *produtos.Produto `json:"produto"`
	Mercado        *mercados.Mercado `json:"mercado"`
	PrecoUnitario  float32           `json:"preco_unitario"`
	NivelConfianca int32             `json:"nivel_confianca"`
	CreatedAt      time.Time         `json:"created_at"`
	ModifiedAt     time.Time         `json:"modified_at"`
	DeletedAt      *time.Time        `json:"deleted_at"`
}
