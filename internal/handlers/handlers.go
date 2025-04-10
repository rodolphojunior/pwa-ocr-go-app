package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"pwaocr/internal/ocr"
	"pwaocr/internal/ai"
	"pwaocr/internal/db"
    "pwaocr/internal/db/models"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10MB
	file, handler, err := r.FormFile("imagem")
	if err != nil {
		http.Error(w, "Erro ao obter o arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	savedPath, err := salvarImagem(file, handler)
	if err != nil {
		http.Error(w, "Erro ao salvar imagem", http.StatusInternalServerError)
		return
	}

	texto, err := ocr.ExtrairTexto(savedPath)
	if err != nil {
		http.Error(w, "Erro no OCR", http.StatusInternalServerError)
		return
	}

	dados, err := ai.ExtrairCampos(texto)
	if err != nil {
		http.Error(w, "Erro ao extrair dados com IA", http.StatusInternalServerError)
		return
	}

	err = db.SalvarNotaFiscal(dados)
	if err != nil {
		http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "sucesso"})
}

func salvarImagem(file multipart.File, handler *multipart.FileHeader) (string, error) {
	dir := "uploads"
	os.MkdirAll(dir, os.ModePerm)
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
	path := filepath.Join(dir, filename)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return "", err
	}

	return path, nil
}

func ListarNotasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	notas := []models.NotaFiscal{}
	db.Conn.Preload("Itens").Find(&notas)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notas)
}

func DeleteNotasHandler(w http.ResponseWriter, r *http.Request) {
	db.Conn.Exec("DELETE FROM item_nota_fiscals") // Nome real da tabela dos itens
	db.Conn.Exec("DELETE FROM nota_fiscals")      // Nome real da tabela das notas
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"Notas apagadas com sucesso"}`))
}



