package mercados

import (
	"strconv"
	"strings"
)

type Mercado struct {
	ID        int64   `json:"id"`
	PlaceID   string  `json:"place_id"`
	Nome      string  `json:"nome"`
	Endereco  string  `json:"endereco"`
	Cidade    string  `json:"cidade"`
	Bairro    string  `json:"bairro"`
	Numero    int     `json:"numero"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    int     `json:"status"`
}

type PlaceGoogle struct {
	ID        string  `json:"id"`
	Nome      string  `json:"nome"`
	Endereco  string  `json:"endereco"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (mkt PlaceGoogle) ParseToMercado() *Mercado {

	splitEndereco := strings.Split(mkt.Endereco, ",")
	rua := splitEndereco[0]
	cidade := splitEndereco[2]

	numeroBairro := strings.Split(splitEndereco[1], "-")
	numero, _ := strconv.Atoi(strings.TrimSpace(numeroBairro[0]))
	bairro := strings.TrimSpace(numeroBairro[1])

	return &Mercado{
		PlaceID:   mkt.ID,
		Nome:      mkt.Nome,
		Endereco:  rua,
		Cidade:    cidade,
		Bairro:    bairro,
		Numero:    numero,
		Latitude:  mkt.Latitude,
		Longitude: mkt.Longitude,
	}
}
