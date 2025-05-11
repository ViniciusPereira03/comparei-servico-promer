package dto

import "main/internal/domain/produtos"

type CreateProductDTO struct {
	Nome       string  `json:"nome" validate:"required"`
	Marca      string  `json:"marca" validate:"required"`
	Quantidade float32 `json:"quantidade" validate:"required"`
	Unidade    string  `json:"unidade" validate:"required"`
	BarCode    string  `json:"bar_code" validate:"required"`
	Latitude   float64 `json:"latitude" validate:"required"`
	Longitude  float64 `json:"longitude" validate:"required"`
	Preco      float32 `json:"preco" validate:"required"`
	Foto       string  `json:"foto" validate:"required"`
}

func (dto *CreateProductDTO) ParseToProduct() *produtos.Produto {
	return &produtos.Produto{
		Nome:       dto.Nome,
		Marca:      dto.Marca,
		Quantidade: dto.Quantidade,
		Unidade:    dto.Unidade,
		BarCode:    dto.BarCode,
		Latitude:   dto.Latitude,
		Longitude:  dto.Longitude,
		Preco:      dto.Preco,
		Foto:       dto.Foto,
	}
}
