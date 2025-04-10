package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
    "strings"
	"os"

	"pwaocr/internal/db/models"
)

func ExtrairCampos(texto string) (*models.NotaFiscalDados, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("variável de ambiente OPENAI_API_KEY não definida")
	}

	// Lê o prompt de um arquivo externo
	rawPrompt, err := os.ReadFile("prompt.txt")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler prompt.txt: %v", err)
	}

	prompt := fmt.Sprintf(string(rawPrompt), texto)

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
	fmt.Println("Resposta bruta da IA:", string(body))

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.Unmarshal(body, &result)

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("IA não retornou nenhuma resposta")
	}

	var dados models.NotaFiscalDados
	rawContent := strings.TrimSpace(result.Choices[0].Message.Content)

	if strings.HasPrefix(rawContent, "```json") {
		rawContent = strings.TrimPrefix(rawContent, "```json")
		rawContent = strings.TrimSuffix(rawContent, "```")
		rawContent = strings.TrimSpace(rawContent)
	}

	err = json.Unmarshal([]byte(rawContent), &dados)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON da IA: %v", err)
	}

	return &dados, nil
}

