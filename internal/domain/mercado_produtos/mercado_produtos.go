package mercadoprodutos

import "time"

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
