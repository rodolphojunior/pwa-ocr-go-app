// internal/db/db.go
package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"pwaocr/internal/db/models"
)

var Conn *gorm.DB

func Connect() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "notas.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}
	Conn = db

	db.AutoMigrate(&models.NotaFiscal{}, &models.ItemNotaFiscal{})
	log.Println("Banco de dados conectado e migrado.")
}

func SalvarNotaFiscal(dados *models.NotaFiscalDados) error {
	nota := models.NotaFiscal{
		Empresa:     dados.Empresa,
		CNPJ:        dados.CNPJ,
		Endereco:    dados.Endereco,
		DataEmissao: dados.DataEmissao,
		ValorTotal:  dados.ValorTotal,
	}
	result := Conn.Create(&nota)
	if result.Error != nil {
		return result.Error
	}

	for _, item := range dados.Itens {
		i := models.ItemNotaFiscal{
			NotaID:       nota.ID,
			Descricao:    item.Descricao,
			Quantidade:   item.Quantidade,
			ValorUnitario: item.ValorUnitario,
			ValorTotal:   item.ValorTotal,
		}
		Conn.Create(&i)
	}

	return nil
}

