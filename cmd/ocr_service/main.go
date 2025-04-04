package main

import (
	"log"
	"net/http"
    "github.com/joho/godotenv"
	apphttp "pwaocr/internal/http"
	"pwaocr/internal/db"
)

func main() {

    // Carrega as variáveis do .env
    if err := godotenv.Load(); err != nil {
        log.Println("Aviso: .env não carregado (pode já estar definido no sistema)")
    }

	db.Connect() // conecta ao banco de dados

	apphttp.HandleRoutes() // define os endpoints

	log.Println("Servidor iniciado em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

