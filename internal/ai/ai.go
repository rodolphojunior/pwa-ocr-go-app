// internal/ai/ai.go
package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
    "pwaocr/internal/db/models"
)

type Item struct {
	Descricao     string  `json:"descricao"`
	Quantidade    int     `json:"quantidade"`
	ValorUnitario float64 `json:"valor_unitario"`
	ValorTotal    float64 `json:"valor_total"`
}


func ExtrairCampos(texto string) (*models.NotaFiscalDados, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("variável de ambiente OPENAI_API_KEY não definida")
	}

	prompt := fmt.Sprintf(`Você é um extrator de dados de notas fiscais brasileiras. Receberá o texto OCR de uma nota fiscal e deve retornar os seguintes campos em JSON: 'empresa', 'cnpj', 'endereco', 'data_emissao', 'itens' (com 'descricao', 'quantidade', 'valor_unitario' e 'valor_total'), e 'valor_total'. Texto OCR: %s`, texto)

	payload := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "system", "content": "Você é um assistente útil."},
			{"role": "user", "content": prompt},
		},
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

    fmt.Println("Resposta bruta da IA:", string(body))

	json.Unmarshal(body, &result)

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("IA não retornou nenhuma resposta")
	}

    var dados models.NotaFiscalDados
    err = json.Unmarshal([]byte(result.Choices[0].Message.Content), &dados)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON da IA: %v", err)
	}

	return &dados, nil
}

