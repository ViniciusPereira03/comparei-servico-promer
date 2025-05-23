package dto

import mercadoprodutos "main/internal/domain/mercado_produtos"

type ConfirmaMercadoProdutoDTO struct {
	IdProduto int64   `json:"nome" validate:"required"`
	IdMercado int64   `json:"marca" validate:"required"`
	Preco     float32 `json:"preco" validate:"required"`
}

func (dto *ConfirmaMercadoProdutoDTO) ParseToMercadoProdutos() *mercadoprodutos.MercadoProdutos {
	return &mercadoprodutos.MercadoProdutos{
		ProdutoID:     dto.IdProduto,
		MercadoID:     dto.IdMercado,
		PrecoUnitario: dto.Preco,
	}
}
