package controllers

import (
	"api/src/autenticacao"
	"api/src/banco"
	"api/src/config"
	"api/src/repositorios"
	"api/src/respostas"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AbacatePayProduct representa um produto para o AbacatePay
type AbacatePayProduct struct {
	ExternalID  string `json:"externalId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"` // Preço em centavos
}

// AbacatePayCustomer representa os dados do cliente
type AbacatePayCustomer struct {
	Name      string `json:"name"`
	Cellphone string `json:"cellphone"`
	Email     string `json:"email"`
	TaxID     string `json:"taxId"`
}

// AbacatePayRequest representa a requisição para criar cobrança
type AbacatePayRequest struct {
	Frequency     string              `json:"frequency"`
	Methods       []string            `json:"methods"`
	Products      []AbacatePayProduct `json:"products"`
	ReturnURL     string              `json:"returnUrl"`
	CompletionURL string              `json:"completionUrl"`
	Customer      AbacatePayCustomer  `json:"customer"`
	AllowCoupons  bool                `json:"allowCoupons"`
	ExternalID    string              `json:"externalId"`
}

// AbacatePayResponse representa a resposta da API do AbacatePay
type AbacatePayResponse struct {
	Error interface{} `json:"error"`
	Data  struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"data"`
}

// GerarIDUnico gera um ID único para usar como externalId
func GerarIDUnico() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("pedido_%d_%s", timestamp, randomHex)
}

// CriarCobrancaAbacatePay cria uma cobrança no AbacatePay e retorna a URL de pagamento
func CriarCobrancaAbacatePay(w http.ResponseWriter, r *http.Request) {
	// Extrair ID do usuário do token
	usuarioID, erro := autenticacao.ExtrairUsuarioID(r)
	if erro != nil {
		respostas.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	// Ler dados da requisição (dados do cliente e endereço)
	corpoRequisicao, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	var dadosCheckout struct {
		Cliente struct {
			Email      string `json:"email"`
			Nome       string `json:"nome"`
			CPF        string `json:"cpf"`
			Telefone   string `json:"telefone"`
			Nascimento string `json:"nascimento"`
		} `json:"cliente"`
		Endereco struct {
			CEP         string `json:"cep"`
			Endereco    string `json:"endereco"`
			Numero      string `json:"numero"`
			Complemento string `json:"complemento"`
			Bairro      string `json:"bairro"`
			Cidade      string `json:"cidade"`
			Estado      string `json:"estado"`
		} `json:"endereco"`
	}

	if erro = json.Unmarshal(corpoRequisicao, &dadosCheckout); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	// Validar CPF (não pode ser todos números iguais)
	cpfLimpo := dadosCheckout.Cliente.CPF
	if cpfLimpo == "" {
		respostas.Erro(w, http.StatusBadRequest, fmt.Errorf("CPF é obrigatório"))
		return
	}

	// Verificar se o CPF não é sequência de números iguais (ex: 11111111111)
	todosIguais := true
	if len(cpfLimpo) > 0 {
		primeiro := cpfLimpo[0]
		for i := 1; i < len(cpfLimpo); i++ {
			if cpfLimpo[i] != primeiro {
				todosIguais = false
				break
			}
		}
	}

	if todosIguais {
		respostas.Erro(w, http.StatusBadRequest, fmt.Errorf("CPF inválido: não pode ser sequência de números iguais (ex: 11111111111)"))
		return
	}

	// Conectar ao banco para buscar carrinho
	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	// Buscar resumo do carrinho
	repositorioCarrinho := repositorios.NovoRepositorioDeCarrinho(db)
	resumo, erro := repositorioCarrinho.BuscarResumoCarrinho(usuarioID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if len(resumo.Itens) == 0 {
		respostas.Erro(w, http.StatusBadRequest, fmt.Errorf("carrinho vazio"))
		return
	}

	// Preparar produtos para o AbacatePay
	products := []AbacatePayProduct{}
	for _, item := range resumo.Itens {
		// Converter preço para centavos (AbacatePay trabalha com centavos)
		precoEmCentavos := int(item.Produto.Preco * 100)

		product := AbacatePayProduct{
			ExternalID:  fmt.Sprintf("prod_%d", item.Produto.ID),
			Name:        item.Produto.Nome,
			Description: fmt.Sprintf("%s - Tamanho: %s", item.Produto.Descricao, item.Produto.Tamanho),
			Quantity:    item.Quantidade,
			Price:       precoEmCentavos,
		}
		products = append(products, product)
	}

	// Adicionar frete se necessário
	if resumo.ValorTotal < 500 {
		products = append(products, AbacatePayProduct{
			ExternalID:  "frete_001",
			Name:        "Frete",
			Description: "Frete de entrega",
			Quantity:    1,
			Price:       3000, // R$ 30,00 em centavos
		})
	}

	// Gerar ID único para o pedido
	externalID := GerarIDUnico()

	// Usar configurações centralizadas
	abacateToken := config.AbacatePayToken
	baseURL := config.BaseURL

	// Preparar requisição para o AbacatePay
	abacateRequest := AbacatePayRequest{
		Frequency:     "ONE_TIME",
		Methods:       []string{"PIX"},
		Products:      products,
		ReturnURL:     fmt.Sprintf("%s/sacola", baseURL),
		CompletionURL: fmt.Sprintf("%s/checkout/confirmacao", baseURL),
		Customer: AbacatePayCustomer{
			Name:      dadosCheckout.Cliente.Nome,
			Cellphone: dadosCheckout.Cliente.Telefone,
			Email:     dadosCheckout.Cliente.Email,
			TaxID:     dadosCheckout.Cliente.CPF,
		},
		AllowCoupons: false,
		ExternalID:   externalID,
	}

	// Converter para JSON
	jsonData, erro := json.Marshal(abacateRequest)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Log para debug
	fmt.Printf("\n=== AbacatePay Request ===\n%s\n========================\n", string(jsonData))

	// Fazer requisição para o AbacatePay
	req, erro := http.NewRequest("POST", "https://api.abacatepay.com/v1/billing/create", bytes.NewBuffer(jsonData))
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", abacateToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, erro := client.Do(req)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer resp.Body.Close()

	// Ler resposta
	respBody, erro := io.ReadAll(resp.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	// Log da resposta do AbacatePay
	fmt.Printf("\n=== AbacatePay Response ===\nStatus: %d\nBody: %s\n========================\n", resp.StatusCode, string(respBody))

	// Verificar se houve erro na API do AbacatePay
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respostas.Erro(w, resp.StatusCode, fmt.Errorf("erro ao criar cobrança: %s", string(respBody)))
		return
	}

	// Parse da resposta
	var abacateResp AbacatePayResponse
	if erro = json.Unmarshal(respBody, &abacateResp); erro != nil {
		fmt.Printf("ERRO ao fazer parse da resposta: %v\n", erro)
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	fmt.Printf("\n=== URL de Pagamento ===\n%s\n========================\n", abacateResp.Data.URL)

	// Retornar URL de pagamento e dados para o frontend
	respostas.JSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"paymentUrl": abacateResp.Data.URL,
		"billingId":  abacateResp.Data.ID,
		"externalId": externalID,
	})
}
