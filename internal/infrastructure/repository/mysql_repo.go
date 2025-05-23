package repository

import (
	"database/sql"
	"fmt"
	mercadoprodutos "main/internal/domain/mercado_produtos"
	"main/internal/domain/mercados"
	"main/internal/domain/produtos"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) CreateProduct(produto *produtos.Produto) (int64, error) {
	result, err := r.db.Exec("INSERT INTO produtos (nome, marca, quantidade, unidade, bar_code) VALUES (?, ?, ?, ?, ?)",
		produto.Nome,
		produto.Marca,
		produto.Quantidade,
		produto.Unidade,
		produto.BarCode,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (r *MySQLRepository) CreateMarketProduct(mercado *mercados.Mercado, produto *produtos.Produto) (int64, error) {
	result, err := r.db.Exec("INSERT INTO mercado_produtos (id_mercado, id_produto, preco_unitario, nivel_confianca) VALUES (?, ?, ?, ?)",
		mercado.ID,
		produto.ID,
		produto.Preco,
		100,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return id, err
}

func (r *MySQLRepository) GetMarketProduct(mercadoId int64, produtoId int64) (*mercadoprodutos.MercadoProdutos, error) {
	mercado_produto := &mercadoprodutos.MercadoProdutos{} // aloca memória

	row := r.db.QueryRow(`
		SELECT id, id_mercado, id_produto, preco_unitario, nivel_confianca, created_at, modified_at, deleted_at
		FROM mercado_produtos WHERE id_mercado = ? AND id_produto = ?
	`, mercadoId, produtoId)

	err := row.Scan(
		&mercado_produto.ID,
		&mercado_produto.MercadoID,
		&mercado_produto.ProdutoID,
		&mercado_produto.PrecoUnitario,
		&mercado_produto.NivelConfianca,
		&mercado_produto.CreatedAt,
		&mercado_produto.ModifiedAt,
		&mercado_produto.DeletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("GetMarketProduct.Scan: %w", err)
	}

	return mercado_produto, nil
}

func (r *MySQLRepository) UpdateMarketProduct(mercado_produtos *mercadoprodutos.MercadoProdutos) error {
	_, err := r.db.Exec("UPDATE mercado_produtos SET preco_unitario = ?, nivel_confianca = ? WHERE id = ?", mercado_produtos.PrecoUnitario, mercado_produtos.NivelConfianca, mercado_produtos.ID)
	return err
}

func (r *MySQLRepository) GetProductByBarcode(barcode string) (*produtos.Produto, error) {
	var p produtos.Produto
	row := r.db.QueryRow(`
		SELECT id, nome, marca, quantidade, unidade, bar_code, created_at, modified_at, deleted_at 
		FROM produtos
		WHERE bar_code = ?
	`, barcode)

	err := row.Scan(
		&p.ID,
		&p.Nome,
		&p.Marca,
		&p.Quantidade,
		&p.Unidade,
		&p.BarCode,
		&p.CreatedAt,
		&p.ModifiedAt,
		&p.DeletedAt,
	)
	if err != nil {
		return &produtos.Produto{}, fmt.Errorf("GetProductByBarcode.Scan: %w", err)
	}
	return &p, nil
}

func (r *MySQLRepository) CreateMarket(mercado *mercados.Mercado) (int64, error) {
	point := fmt.Sprintf("POINT(%f %f)", mercado.Longitude, mercado.Latitude)

	result, err := r.db.Exec("INSERT INTO mercados (nome, endereco, cidade, bairro, numero, latitude, longitude, local) VALUES (?, ?, ?, ?, ?, ?, ?, ST_PointFromText(?))",
		mercado.Nome,
		mercado.Endereco,
		mercado.Cidade,
		mercado.Bairro,
		mercado.Numero,
		mercado.Latitude,
		mercado.Longitude,
		point,
	)
	if err != nil {
		fmt.Println("ERRO: ", err)
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

func (r *MySQLRepository) GetMarketByCoordinates(lat float64, lng float64) (mercados.PlaceGoogle, error) {
	return mercados.PlaceGoogle{}, nil
}

func (r *MySQLRepository) SearchMarketByCoordinates(lat float64, lng float64) (*mercados.Mercado, error) {
	var m mercados.Mercado
	row := r.db.QueryRow(`
		SELECT id, nome, endereco, cidade, bairro, numero, ST_X(local) as latitude, ST_Y(local) as longitude, status
		FROM mercados
		WHERE ST_Distance_Sphere(local, POINT(?, ?)) <= 50
	`, lng, lat)

	err := row.Scan(
		&m.ID,
		&m.Nome,
		&m.Endereco,
		&m.Cidade,
		&m.Bairro,
		&m.Numero,
		&m.Latitude,
		&m.Longitude,
		&m.Status,
	)
	if err != nil {
		return &mercados.Mercado{}, fmt.Errorf("GetMarketByCoordinates.Scan: %w", err)
	}
	return &m, nil
}

func (r *MySQLRepository) IdetificarProduto(produto *produtos.ProdutoFoto) (*produtos.Produto, error) {
	return &produtos.Produto{}, nil
}

func (r *MySQLRepository) ConfirmarValor(data *mercadoprodutos.MercadoProdutos, userId string) (int64, error) {

	mercado_produto := &mercadoprodutos.MercadoProdutos{} // aloca memória

	row := r.db.QueryRow(`
		SELECT id, id_mercado, id_produto, preco_unitario, nivel_confianca, created_at, modified_at, deleted_at
		FROM mercado_produtos WHERE id_mercado = ? AND id_produto = ?
	`, data.MercadoID, data.ProdutoID)

	err := row.Scan(
		&mercado_produto.ID,
		&mercado_produto.MercadoID,
		&mercado_produto.ProdutoID,
		&mercado_produto.PrecoUnitario,
		&mercado_produto.NivelConfianca,
		&mercado_produto.CreatedAt,
		&mercado_produto.ModifiedAt,
		&mercado_produto.DeletedAt,
	)

	if err != nil {
		return 0, err
	}
	novoNivel := mercado_produto.NivelConfianca + 10
	if novoNivel > 100 {
		novoNivel = 100
	}

	if mercado_produto.PrecoUnitario != data.PrecoUnitario {
		_, err := r.db.Exec("UPDATE mercado_produtos SET preco_unitario = ?, nivel_confianca = ? WHERE id = ?", data.PrecoUnitario, novoNivel, mercado_produto.ID)
		return mercado_produto.ID, err
	} else {
		_, err := r.db.Exec("UPDATE mercado_produtos SET nivel_confianca = ? WHERE id = ?", novoNivel, mercado_produto.ID)
		return mercado_produto.ID, err
	}
}
