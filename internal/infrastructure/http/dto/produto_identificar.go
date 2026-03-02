package dto

import "main/internal/domain/produtos"

type ProdutoFotoDTO struct {
	Foto string `json:"foto" validate:"required"`
}

func (dto *ProdutoFotoDTO) ParseToPhoto() *produtos.ProdutoFoto {
	return &produtos.ProdutoFoto{
		Foto: dto.Foto,
	}
}
