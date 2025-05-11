package produtos_interface

import (
	mercadoprodutos "main/internal/domain/mercado_produtos"
	"main/internal/domain/mercados"
	"main/internal/domain/produtos"
)

type ProdutosRepository interface {
	CreateProduct(produto *produtos.Produto) (int64, error)
	CreateMarketProduct(mercado *mercados.Mercado, produto *produtos.Produto) (int64, error)
	GetMarketProduct(mercadoId int64, produtoId int64) (*mercadoprodutos.MercadoProdutos, error)
	UpdateMarketProduct(mercado_produtos *mercadoprodutos.MercadoProdutos) error
	GetProductByBarcode(barcode string) (*produtos.Produto, error)
}
