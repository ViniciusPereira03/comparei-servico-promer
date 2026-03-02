package produtos

import "time"

type Produto struct {
	ID         int64      `json:"id"`
	Nome       string     `json:"nome"`
	Marca      string     `json:"marca"`
	Quantidade float32    `json:"quantidade"`
	Unidade    string     `json:"unidade"`
	BarCode    string     `json:"bar_code"`
	Latitude   float64    `json:"latitude"`
	Longitude  float64    `json:"longitude"`
	Preco      float32    `json:"preco"`
	Foto       string     `json:"foto"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt time.Time  `json:"modified_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type ProdutoFoto struct {
	Foto string `json:"foto"`
}
