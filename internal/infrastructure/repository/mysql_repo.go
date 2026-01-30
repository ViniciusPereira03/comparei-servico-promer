package repository

import (
	"database/sql"
	"fmt"
	mercadoprodutos "main/internal/domain/mercado_produtos"
	"main/internal/domain/mercados"
	"main/internal/domain/produtos"
	"main/internal/domain/user"
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

func (r *MySQLRepository) GetMarketProductId(mercadoProdutoId int64) (*mercadoprodutos.MercadoProdutos, error) {
	mercado_produto := &mercadoprodutos.MercadoProdutos{} // aloca memória

	row := r.db.QueryRow(`
		SELECT id, id_mercado, id_produto, preco_unitario, nivel_confianca, created_at, modified_at, deleted_at
		FROM mercado_produtos WHERE id = ?
	`, mercadoProdutoId)

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

func (r *MySQLRepository) SearchProductsByText(text string) ([]produtos.Produto, error) {
	var produtosList []produtos.Produto

	// Monta o pattern para o LIKE
	likePattern := "%" + text + "%"

	rows, err := r.db.Query(`
		SELECT id, nome, marca, quantidade, unidade, bar_code, created_at, modified_at, deleted_at
		FROM produtos
		WHERE nome LIKE ? OR marca LIKE ?
	`, likePattern, likePattern)
	if err != nil {
		return nil, fmt.Errorf("SearchProductsByText.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p produtos.Produto
		err := rows.Scan(
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
			return nil, fmt.Errorf("SearchProductsByText.Scan: %w", err)
		}
		produtosList = append(produtosList, p)
	}

	// Checa se houve algum erro durante a iteração
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("SearchProductsByText.rows.Err: %w", err)
	}

	return produtosList, nil
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

func (r *MySQLRepository) GetMarketByCoordinates(lat float64, lng float64, radius int) (mercados.PlaceGoogle, error) {
	return mercados.PlaceGoogle{}, nil
}

func (r *MySQLRepository) SearchMarketByCoordinates(lat float64, lng float64, radius int) (*mercados.Mercado, error) {
	var m mercados.Mercado
	row := r.db.QueryRow(`
		SELECT id, nome, endereco, cidade, bairro, numero, ST_X(local) as latitude, ST_Y(local) as longitude, status
		FROM mercados
		WHERE ST_Distance_Sphere(local, POINT(?, ?)) <= ?
	`, lng, lat, radius)

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

func (r *MySQLRepository) SearchMarketsByCoordinates(
	lat float64,
	lng float64,
	radius int,
) ([]mercados.Mercado, error) {
	fmt.Println(fmt.Sprintf("BUSCANDO MERCADOS EM %v %v %vm", lat, lng, radius))
	var mercadosEncontrados []mercados.Mercado

	rows, err := r.db.Query(`
		SELECT 
			id,
			nome,
			endereco,
			cidade,
			bairro,
			numero,
			ST_X(local) AS latitude,
			ST_Y(local) AS longitude,
			status
		FROM mercados
		WHERE ST_Distance_Sphere(local, POINT(?, ?)) <= ?
	`, lng, lat, radius)
	if err != nil {
		return nil, fmt.Errorf("SearchMarketByCoordinates.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m mercados.Mercado

		err := rows.Scan(
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
			return nil, fmt.Errorf("SearchMarketByCoordinates.Scan: %w", err)
		}

		mercadosEncontrados = append(mercadosEncontrados, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("SearchMarketByCoordinates.Rows: %w", err)
	}

	fmt.Println("mercadosEncontrados: ", mercadosEncontrados)
	return mercadosEncontrados, nil
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

func (r *MySQLRepository) GetMarketsByProduct(produto *produtos.Produto) ([]*mercadoprodutos.MercadoProdutosCompleto, error) {
	var mercadosList []*mercadoprodutos.MercadoProdutosCompleto

	rows, err := r.db.Query(`
		SELECT m.id, m.nome, m.endereco, m.cidade, m.bairro, m.numero, m.latitude, m.longitude, m.status,
			   mp.id as id_mercado_produtos, mp.preco_unitario, mp.nivel_confianca, mp.created_at, mp.modified_at
		FROM mercado_produtos mp
		INNER JOIN mercados m ON m.id = mp.id_mercado
		WHERE mp.id_produto = ?
	`, produto.ID)
	if err != nil {
		return nil, fmt.Errorf("GetMarketsByProduct.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m mercados.Mercado
		var mp mercadoprodutos.MercadoProdutosCompleto

		err := rows.Scan(
			&m.ID,
			&m.Nome,
			&m.Endereco,
			&m.Cidade,
			&m.Bairro,
			&m.Numero,
			&m.Latitude,
			&m.Longitude,
			&m.Status,
			&mp.ID,
			&mp.PrecoUnitario,
			&mp.NivelConfianca,
			&mp.CreatedAt,
			&mp.ModifiedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("GetMarketsByProduct.Scan: %w", err)
		}

		mp.Produto = produto
		mp.Mercado = &m

		mercadosList = append(mercadosList, &mp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("GetMarketsByProduct.RowsErr: %w", err)
	}

	return mercadosList, nil
}

func (r *MySQLRepository) GetMarketByID(marketID int64) (*mercados.Mercado, error) {
	var m mercados.Mercado
	row := r.db.QueryRow(`
		SELECT id, nome, endereco, cidade, bairro, numero, ST_X(local) as latitude, ST_Y(local) as longitude, status
		FROM mercados
		WHERE id = ?
	`, marketID)

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
		return &mercados.Mercado{}, fmt.Errorf("GetMarketByID.Scan: %w", err)
	}
	return &m, nil
}

// Users
func (r *MySQLRepository) CreateUser(user *user.User) error {
	_, err := r.db.Exec("INSERT INTO users (id, ray_distance, status) VALUES (?, ?, ?)", user.ID, user.RayDistance, user.Status)
	return err
}

func (r *MySQLRepository) EditUser(user *user.User) error {
	_, err := r.db.Exec("UPDATE users SET ray_distance = ?, status = ? WHERE id = ?", user.RayDistance, user.Status, user.ID)
	return err
}

func (r *MySQLRepository) GetUser(id string) (*user.User, error) {
	user := &user.User{} // aloca memória

	row := r.db.QueryRow(`
		SELECT id, ray_distance, status
		FROM users WHERE id = ? and deleted_at IS NULL
	`, id)

	err := row.Scan(
		&user.ID,
		&user.RayDistance,
		&user.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("GetUser.Scan: %w", err)
	}

	return user, nil
}
