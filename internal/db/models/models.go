// internal/db/models/models.go
package models

import "gorm.io/gorm"

type NotaFiscal struct {
	gorm.Model
	Empresa     string
	CNPJ        string
	Endereco    string
	DataEmissao string
	ValorTotal  float64
	Itens       []ItemNotaFiscal `gorm:"foreignKey:NotaID"`
}

type ItemNotaFiscal struct {
	gorm.Model
	NotaID       uint
	Descricao    string
	Quantidade   int
	ValorUnitario float64
	ValorTotal   float64
}

type Item struct {
	Descricao     string  `json:"descricao"`
	Quantidade    int     `json:"quantidade"`
	ValorUnitario float64 `json:"valor_unitario"`
	ValorTotal    float64 `json:"valor_total"`
}

type NotaFiscalDados struct {
	Empresa     string  `json:"empresa"`
	CNPJ        string  `json:"cnpj"`
	Endereco    string  `json:"endereco"`
	DataEmissao string  `json:"data_emissao"`
	Itens       []Item  `json:"itens"`
	ValorTotal  float64 `json:"valor_total"`
}

