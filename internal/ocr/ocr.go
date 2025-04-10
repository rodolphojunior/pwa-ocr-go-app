package ocr

import (
	"fmt"
	"github.com/otiai10/gosseract/v2"
)

func ExtrairTexto(caminhoImagem string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetLanguage("por")
	client.SetImage(caminhoImagem)

	texto, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("erro ao realizar OCR: %v", err)
	}

	return texto, nil
}

