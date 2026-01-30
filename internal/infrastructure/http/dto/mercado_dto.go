package dto

import "main/internal/domain/mercados"

type MercadoDTO struct {
	Nome      string  `json:"nome"`
	Endereco  string  `json:"endereco"`
	Cidade    string  `json:"cidade"`
	Bairro    string  `json:"bairro"`
	Numero    int     `json:"numero"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (dto *MercadoDTO) ParseToMarket() *mercados.Mercado {
	return &mercados.Mercado{
		Nome:      dto.Nome,
		Endereco:  dto.Endereco,
		Cidade:    dto.Cidade,
		Bairro:    dto.Bairro,
		Numero:    dto.Numero,
		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,
	}
}
