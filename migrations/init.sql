CREATE DATABASE IF NOT EXISTS promerdb;

-- Tabela de produtos
CREATE TABLE IF NOT EXISTS produtos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    marca VARCHAR(255) NOT NULL,
    quantidade DECIMAL(10,2) NOT NULL,
    unidade VARCHAR(100) NOT NULL,
    bar_code VARCHAR(255) UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
);

-- Tabela de mercados
CREATE TABLE IF NOT EXISTS mercados (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    endereco VARCHAR(255) NOT NULL,
    cidade VARCHAR(255) NOT NULL,
    bairro VARCHAR(255) NOT NULL,
    numero INT NOT NULL DEFAULT 0,
    latitude VARCHAR(255) NOT NULL,
    longitude VARCHAR(255) NOT NULL,
    local POINT NOT NULL,
    status INT DEFAULT 1 NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    SPATIAL INDEX idx_local (local)
);

-- Tabela de produtos do mercado
CREATE TABLE IF NOT EXISTS mercado_produtos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_mercado INT NOT NULL,
    id_produto INT NOT NULL,
    preco_unitario DOUBLE(10, 2) NOT NULL,
    nivel_confianca INT DEFAULT 100,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (id_mercado) REFERENCES mercados(id),
    FOREIGN KEY (id_produto) REFERENCES produtos(id)
);
