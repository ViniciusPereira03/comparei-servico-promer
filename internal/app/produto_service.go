package app

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	mercadoprodutos "main/internal/domain/mercado_produtos"
	"main/internal/domain/mercados"
	mercados_interface "main/internal/domain/mercados/interface"
	"main/internal/domain/produtos"
	produtos_interface "main/internal/domain/produtos/interface"
	"main/internal/infrastructure/messaging/publisher"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
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

func (s *ProdutosService) CreateProduct(product *produtos.Produto, userId string) (int64, error) {

	radius := 100
	mercado, err := s.mercadoService.SearchMarketByCoordinates(product.Latitude, product.Longitude, radius)
	if err != nil {
		mkt, err := s.mercadoService.GetMarketByCoordinates(product.Latitude, product.Longitude, radius)
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

		mercadoProduto, err := s.GetMarketProductId(idMercadoProduto)
		if err != nil {
			return 0, err
		}
		err_pub := publisher.PubNewProduct(mercadoProduto, userId)
		if err_pub != nil {
			log.Println("[ERRO PUB] ", err_pub)
		}
	} else {
		idMercadoProduto = mercado_produto.ID
		mercado_produto.PrecoUnitario = product.Preco
		mercado_produto.NivelConfianca = 100

		err = s.UpdateMarketProduct(mercado_produto)
		if err != nil {
			return 0, err
		}

		err_pub := publisher.PubUpdateProduct(idMercadoProduto, userId)
		if err_pub != nil {
			log.Println("[ERRO PUB] ", err_pub)
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

func (s *ProdutosService) GetMarketProductId(mercadoProdutoId int64) (*mercadoprodutos.MercadoProdutos, error) {
	return s.mysqlRepo.GetMarketProductId(mercadoProdutoId)
}

func (s *ProdutosService) UpdateMarketProduct(mercado_produto *mercadoprodutos.MercadoProdutos) error {
	return s.mysqlRepo.UpdateMarketProduct(mercado_produto)
}

func (s *ProdutosService) IdetificarProduto(produto *produtos.ProdutoFoto) (*produtos.Produto, error) {
	caminho, errImg := SalvarImagemBase64(produto.Foto)
	if errImg != nil {
		fmt.Println("ERRO: ", errImg.Error())
		return &produtos.Produto{}, errImg
	}

	// Structs para armazenar a resposta do modelo
	type Content struct {
		Parts []string `json:"Parts"` // Mantendo como string para conversão futura
		Role  string   `json:"Role"`
	}

	type Candidate struct {
		Index            int     `json:"Index"`
		Content          Content `json:"Content"`
		FinishReason     int     `json:"FinishReason"`
		SafetyRatings    any     `json:"SafetyRatings"`
		CitationMetadata any     `json:"CitationMetadata"`
		TokenCount       int     `json:"TokenCount"`
	}

	err := godotenv.Load()
	if err != nil {
		// log.Fatal("Erro ao carregar o arquivo .env")
		return &produtos.Produto{}, err
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		// log.Fatal(err)
		return &produtos.Produto{}, err
	}
	defer client.Close()

	// model := client.GenerativeModel("gemini-1.5-flash")
	model := client.GenerativeModel("gemini-2.5-flash-lite")

	imgData1, err := os.ReadFile(caminho)
	if err != nil {
		// log.Fatal(err)
		return &produtos.Produto{}, err
	}

	prompt := []genai.Part{
		genai.ImageData("jpeg", imgData1),
		// genai.Text("Identify the product on the label with the price and return a JSON object with the product information on the market shelf following the following example: { product: 'Arroz', brand: 'Camil', amount: 5, unity: 'Kg', price: 23.99} If necessary, consider only the retail price"),
		genai.Text("Examine the label and identify the product along with its price. Based on this information, generate a JSON object containing the product details as displayed on the market shelf. The JSON object must follow exactly the structure below:\n\n{\"nome\":\"Arroz\",\n \"marca\":\"Camil\",\n \"quantidade\":5,\n \"unidade\":\"Kg\",\n \"preco\":23.99\n}\n\nEnsure that the measurement unit is always provided as an abbreviation (e.g., Kg, g, L, mL) and never as a full word (e.g., KILO). If multiple prices are available, consider only the retail price. If the retail price is not identified, select the highest price among the identified prices. Additionally, ensure that product names are always in Brazilian Portuguese. Return only the JSON object without any additional text."),
	}
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		// log.Fatal(err)
		return &produtos.Produto{}, err
	}

	if len(resp.Candidates) == 0 {
		// log.Fatal("Nenhuma resposta gerada pelo modelo")
		return &produtos.Produto{}, err
	}

	// Converte a resposta para a struct Candidate
	var firstCandidate Candidate
	jsonData, err := json.Marshal(resp.Candidates[0])
	if err != nil {
		// log.Fatal("Erro ao converter para JSON:", err)
		return &produtos.Produto{}, err
	}

	err = json.Unmarshal(jsonData, &firstCandidate)
	if err != nil {
		// log.Fatal("Erro ao decodificar JSON:", err)
		return &produtos.Produto{}, err
	}

	// Extrai JSON de dentro de Parts
	cleanJSON, err := extractJSONFromParts(firstCandidate.Content.Parts)
	if err != nil {
		// log.Fatal("Erro ao extrair JSON de Parts:", err)
		return &produtos.Produto{}, err
	}

	fmt.Println("Conteúdo JSON formatado:")
	fmt.Println(cleanJSON)

	var product *produtos.Produto
	errJson := json.Unmarshal([]byte(cleanJSON), &product)
	if errJson != nil {
		log.Printf("Erro ao fazer unmarshal do JSON: %v", errJson)
		return &produtos.Produto{}, errJson
	}

	return product, nil
}

func SalvarImagemBase64(base64String string) (string, error) {
	const diretorio = "./tmp"
	nomeArquivo := uuid.New().String()

	var img image.Image
	var data []byte
	var err error

	if strings.HasPrefix(base64String, "data:image/png;base64,") {
		base64String = strings.TrimPrefix(base64String, "data:image/png;base64,")
		data, err = base64.StdEncoding.DecodeString(base64String)
		if err != nil {
			return "", fmt.Errorf("erro ao decodificar base64 PNG: %w", err)
		}
		img, err = png.Decode(bytes.NewReader(data))
		if err != nil {
			return "", fmt.Errorf("erro ao decodificar a imagem PNG: %w", err)
		}
		nomeArquivo += ".png"
	} else if strings.HasPrefix(base64String, "data:image/jpeg;base64,") || strings.HasPrefix(base64String, "data:image/jpg;base64,") {
		base64String = strings.TrimPrefix(base64String, "data:image/jpeg;base64,")
		base64String = strings.TrimPrefix(base64String, "data:image/jpg;base64,")
		data, err = base64.StdEncoding.DecodeString(base64String)
		if err != nil {
			return "", fmt.Errorf("erro ao decodificar base64 JPEG: %w", err)
		}
		img, err = jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return "", fmt.Errorf("erro ao decodificar a imagem JPEG: %w", err)
		}
		nomeArquivo += ".jpg"
	} else {
		return "", fmt.Errorf("formato de imagem não suportado ou prefixo inválido")
	}

	// Cria o diretório ./tmp se não existir
	err = os.MkdirAll(diretorio, 0755)
	if err != nil {
		return "", fmt.Errorf("erro ao criar o diretório ./tmp: %w", err)
	}

	caminhoArquivo := filepath.Join(diretorio, nomeArquivo)

	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return "", fmt.Errorf("erro ao criar o arquivo: %w", err)
	}
	defer arquivo.Close()

	if strings.HasSuffix(caminhoArquivo, ".png") {
		err = png.Encode(arquivo, img)
		if err != nil {
			return "", fmt.Errorf("erro ao salvar a imagem PNG: %w", err)
		}
	} else if strings.HasSuffix(caminhoArquivo, ".jpg") {
		err = jpeg.Encode(arquivo, img, nil)
		if err != nil {
			return "", fmt.Errorf("erro ao salvar a imagem JPEG: %w", err)
		}
	}

	return caminhoArquivo, nil
}

func extractJSONFromParts(parts []string) (string, error) {
	if len(parts) == 0 {
		return "", fmt.Errorf("nenhum conteúdo encontrado em Parts")
	}

	re := regexp.MustCompile("(?s)```json\\n(.*)\\n```")
	matches := re.FindStringSubmatch(parts[0])
	if len(matches) > 1 {
		return matches[1], nil
	}
	return parts[0], nil
}

func (s *ProdutosService) SearchProductsByBarcode(barcode string) ([]*mercadoprodutos.MercadoProdutosCompleto, error) {

	produto, err := s.mysqlRepo.GetProductByBarcode(barcode)
	if err != nil {
		return []*mercadoprodutos.MercadoProdutosCompleto{}, fmt.Errorf("Produto não encontrado")
	}

	mercados, err := s.mysqlRepo.GetMarketsByProduct(produto)

	return mercados, nil
}

func (s *ProdutosService) SearchProductsByText(text string) ([]*mercadoprodutos.MercadoProdutosCompleto, error) {

	produtos, err := s.mysqlRepo.SearchProductsByText(text)
	if err != nil {
		return []*mercadoprodutos.MercadoProdutosCompleto{}, fmt.Errorf("Produto não encontrado")
	}

	var mercados []*mercadoprodutos.MercadoProdutosCompleto

	for _, produto := range produtos {
		mkts, _ := s.mysqlRepo.GetMarketsByProduct(&produto)
		mercados = append(mercados, mkts...)
	}

	return mercados, nil
}
