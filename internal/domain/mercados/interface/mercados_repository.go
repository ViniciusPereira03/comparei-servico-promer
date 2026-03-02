package mercados_interface

import (
	mercadoprodutos "main/internal/domain/mercado_produtos"
	"main/internal/domain/mercados"
)

type MercadosRepository interface {
	CreateMarket(mercado *mercados.Mercado) (int64, error)
	GetMarketByID(marketID int64) (*mercados.Mercado, error)
	GetMarketByCoordinates(lat float64, lng float64, radius int) (mercados.PlaceGoogle, error)
	SearchMarketByCoordinates(lat float64, lng float64, radius int) (*mercados.Mercado, error)
	SearchMarketsByCoordinates(lat float64, lng float64, radius int) ([]mercados.Mercado, error)
	GetMarketProductId(mercadoProdutoId int64) (*mercadoprodutos.MercadoProdutos, error)
	ConfirmarValor(data *mercadoprodutos.MercadoProdutos, userId string) (int64, error)
}
