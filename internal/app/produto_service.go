package app

import (
	mercadoprodutos "main/internal/domain/mercado_produtos"
	"main/internal/domain/mercados"
	mercados_interface "main/internal/domain/mercados/interface"
	"main/internal/domain/produtos"
	produtos_interface "main/internal/domain/produtos/interface"
)

type ProdutosService struct {
	mysqlRepo      produtos_interface.ProdutosRepository
	mercadoService mercados_interface.MercadosRepository
	mongoRepo      produtos_interface.MongoRepository
}

func NewProductsService(
	mysqlRepo produtos_interface.ProdutosRepository,
	mercado mercados_interface.MercadosRepository,
	mongoRepo produtos_interface.MongoRepository,
) *ProdutosService {
	return &ProdutosService{
		mysqlRepo:      mysqlRepo,
		mercadoService: mercado,
		mongoRepo:      mongoRepo,
	}
}

func (s *ProdutosService) CreateProduct(product *produtos.Produto) (int64, error) {

	mercado, err := s.mercadoService.SearchMarketByCoordinates(product.Latitude, product.Longitude)
	if err != nil {
		mkt, err := s.mercadoService.GetMarketByCoordinates(product.Latitude, product.Longitude)
		if err != nil {
			return 0, err
		}

		mercado = mkt.ParseToMercado()

		mercadoId, err := s.mercadoService.CreateMarket(mercado)
		if err != nil {
			return 0, err
		}

		mercado.ID = mercadoId
	}

	var produtoID int64
	produto, err := s.mysqlRepo.GetProductByBarcode(product.BarCode)
	if err != nil {
		produtoID, err = s.mysqlRepo.CreateProduct(product)
		if err != nil {
			return 0, err
		}
		product.ID = produtoID
	} else {
		product.ID = produto.ID
	}

	var idMercadoProduto int64
	mercado_produto, err := s.GetMarketProduct(mercado.ID, product.ID)
	if err != nil {
		idMercadoProduto, err = s.CreateMarketProduct(mercado, product)
		if err != nil {
			return 0, err
		}
	} else {
		mercado_produto.PrecoUnitario = product.Preco
		mercado_produto.NivelConfianca = 100

		err = s.UpdateMarketProduct(mercado_produto)
		if err != nil {
			return 0, err
		}
	}

	err = s.mongoRepo.SaveImage(product.Foto, product.BarCode)
	if err != nil {
		return 0, err
	}

	return idMercadoProduto, nil
}

func (s *ProdutosService) CreateMarketProduct(mercado *mercados.Mercado, produto *produtos.Produto) (int64, error) {
	return s.mysqlRepo.CreateMarketProduct(mercado, produto)
}

func (s *ProdutosService) GetMarketProduct(mercadoId int64, produtoId int64) (*mercadoprodutos.MercadoProdutos, error) {
	return s.mysqlRepo.GetMarketProduct(mercadoId, produtoId)
}

func (s *ProdutosService) UpdateMarketProduct(mercado_produto *mercadoprodutos.MercadoProdutos) error {
	return s.mysqlRepo.UpdateMarketProduct(mercado_produto)
}
