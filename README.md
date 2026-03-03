# Comparei - Serviço Promer 🛒

O **Comparei - Serviço Promer** é um microsserviço de backend responsável pelo gerenciamento de mercados, produtos e usuários. Ele faz parte do ecossistema "Comparei", fornecendo uma API REST robusta e processamento assíncrono de mensagens.

## 🛠️ Tecnologias Utilizadas

Este serviço foi construído utilizando as seguintes tecnologias e ferramentas:

* **Linguagem:** [Go 1.22](https://golang.org/)
* **Bancos de Dados:**
  * **MySQL 8:** Para armazenamento relacional de entidades estruturadas (usuários, mercados, etc).
  * **MongoDB:** Para armazenamento flexível de documentos baseados em coleções (ex: catálogos complexos de produtos).
* **Mensageria e Cache:** **Redis 7** (utilizado para arquitetura orientada a eventos via Pub/Sub).
* **Infraestrutura:** Docker e Docker Compose.
* **Roteamento HTTP:** [Gorilla Mux](https://github.com/gorilla/mux)
* **Integrações Externas:**
  * Google Generative AI (Processamento inteligente de dados de produtos)
  * Google Maps API (Geolocalização para mercados)
  * Autenticação via JWT (JSON Web Tokens)

## ⚙️ Arquitetura do Sistema

A aplicação inicializa dois fluxos principais simultaneamente:
1. **Servidor HTTP:** Expõe os *endpoints* da API REST para lidar com as requisições de Produtos, Mercados e Usuários.
2. **Subscriber (Mensageria):** Uma rotina em *background* (*goroutine*) que escuta eventos do Redis, permitindo comunicação assíncrona com outros microsserviços da rede `comparei_net`.

## 🚀 Como Executar o Projeto Localmente

### Pré-requisitos
* [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados.
* [Go 1.22+](https://golang.org/dl/) (caso queira rodar fora dos contêineres).

### Passo a Passo

1. **Clone o repositório e acesse a pasta:**
```bash
   git clone https://github.com/ViniciusPereira03/comparei-servico-promer
   cd comparei-servico-promer

```

2. **Configuração de Variáveis de Ambiente:**
Crie um arquivo `.env` na raiz do projeto contendo as variáveis necessárias. A aplicação espera, no mínimo, as seguintes chaves (verifique o arquivo de exemplo, se houver):
```env
PORT=8082

# MySQL
MYSQL_HOST=db:3306
MYSQL_USER=root
MYSQL_PASSWORD=root
MYSQL_DB=promerdb

# MongoDB
MONGO_URI=mongodb://root:promerdb@mongo-promer:27017
MONGO_DB_NAME=promerdb
MONGO_COLLECTION=produtos

# Redis
REDIS_MESSAGING_HOST=redis
REDIS_MESSAGING_PORT=6379

```


3. **Executar a Aplicação com `run.sh`:**
Dá permissão de execução ao script (caso ainda não tenha dado) e executa-o. Este ficheiro já está configurado para inicializar a aplicação corretamente com os devidos parâmetros:
```bash
chmod +x run.sh
./run.sh

```

> **Nota:** O script aguardará automaticamente (`wait-for-it.sh`) os bancos de dados estarem prontos antes de iniciar o servidor Go.

4. **Verificar os Logs:**
Após a execução do script, deves ver as seguintes mensagens no terminal confirmando que os serviços estão a rodar:
* `📡 Inicializando subscriber...`
* `🚀 Servidor rodando na porta 8082`


## 📂 Estrutura de Diretórios (Resumo)

* `/config`: Carregamento e validação das variáveis de ambiente (`.env`).
* `/internal`: Código privado da aplicação.
* `/app`: Lógica de negócio e serviços (`produto_service.go`, `mercado_service.go`, etc).
* `/infrastructure`: Camada de adaptação técnica.
* `/http`: Roteadores, *handlers* e DTOs.
* `/messaging`: *Publishers* e *Subscribers* do Redis.
* `/repository`: Implementações do banco de dados (MySQL e MongoDB).
* `/migrations`: Scripts de inicialização do banco de dados relacional (`init.sql`).
* `/tmp`: Armazenamento temporário, como por exemplo, imagens e uploads.
