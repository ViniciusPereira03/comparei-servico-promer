package app

import (
	"context"
	"fmt"
	"log"
	mercadoprodutos "main/internal/domain/mercado_produtos"
	"main/internal/domain/mercados"
	mercados_interface "main/internal/domain/mercados/interface"
	"main/internal/infrastructure/messaging/publisher"
	"os"

	"googlemaps.github.io/maps"
)

type MercadoService struct {
	mysqlRepo mercados_interface.MercadosRepository
}

func NewMercadoService(mysqlRepo mercados_interface.MercadosRepository) *MercadoService {
	return &MercadoService{mysqlRepo: mysqlRepo}
}

func (s *MercadoService) CreateMarket(mercado *mercados.Mercado) (int64, error) {
	return s.mysqlRepo.CreateMarket(mercado)
}

func (s *MercadoService) GetMarketByCoordinates(lat float64, lng float64) (mercados.PlaceGoogle, error) {
	ctx := context.Background()

	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		log.Fatal("defina a variável de ambiente GOOGLE_MAPS_API_KEY")
	}
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("erro ao criar cliente: %v", err)
	}

	var mkt mercados.PlaceGoogle

	radius := 100
	req := &maps.NearbySearchRequest{
		Location: &maps.LatLng{Lat: lat, Lng: lng},

		Radius: uint(radius),
		Type:   "supermarket",
	}

	nearbyResp, err := client.NearbySearch(ctx, req)
	if err != nil {
		log.Printf("PlacesNearby falhou: %v", err)
	}

	if len(nearbyResp.Results) == 0 {
		fmt.Println("Nenhum estabelecimento encontrado via Places API.")
		return mkt, nil
	}

	place := nearbyResp.Results[0]

	placeDetailsReq := &maps.PlaceDetailsRequest{
		PlaceID: place.PlaceID,
		Fields:  []maps.PlaceDetailsFieldMask{maps.PlaceDetailsFieldMaskGeometryLocation},
	}
	placeDetailsResp, err := client.PlaceDetails(ctx, placeDetailsReq)
	if err != nil {
		log.Printf("PlaceDetails falhou: %v", err)
		return mkt, nil
	}
	centralLat := placeDetailsResp.Geometry.Location.Lat
	centralLng := placeDetailsResp.Geometry.Location.Lng

	mkt.ID = place.PlaceID
	mkt.Endereco = place.Vicinity
	mkt.Nome = place.Name
	mkt.Latitude = centralLat
	mkt.Longitude = centralLng

	return mkt, nil
}

func (s *MercadoService) SearchMarketByCoordinates(lat float64, lng float64) (*mercados.Mercado, error) {
	return s.mysqlRepo.SearchMarketByCoordinates(lat, lng)
}

func (s *MercadoService) ConfirmarValor(data *mercadoprodutos.MercadoProdutos, userId string) (int64, error) {
	idMercadoProduto, err := s.mysqlRepo.ConfirmarValor(data, userId)

	if err == nil && idMercadoProduto > 0 {
		err_pub := publisher.PubNewProduct(idMercadoProduto, userId)
		if err_pub != nil {
			log.Println("[ERRO PUB] ", err_pub)
		}
	}

	return idMercadoProduto, err
}
