package mercados_interface

import "main/internal/domain/mercados"

type MercadosRepository interface {
	CreateMarket(mercado *mercados.Mercado) (int64, error)
	GetMarketByCoordinates(lat float64, lng float64) (mercados.PlaceGoogle, error)
	SearchMarketByCoordinates(lat float64, lng float64) (*mercados.Mercado, error)
}
