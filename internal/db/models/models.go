package models

import "gorm.io/gorm"

type NotaFiscal struct {
	gorm.Model
	Empresa     string            `json:"empresa"`
	CNPJ        string            `json:"cnpj"`
	Endereco    string            `json:"endereco"`
	DataEmissao string            `json:"data_emissao"`
	ValorTotal  float64           `json:"valor_total"`
    Itens []ItemNotaFiscal `gorm:"foreignKey:NotaID" json:"itens"`
}

type ItemNotaFiscal struct {
	gorm.Model
	NotaID          uint    `json:"nota_id"`
	Descricao       string  `json:"descricao"`
	Quantidade      int     `json:"quantidade"`
	ValorUnitario   float64 `json:"valor_unitario"`
	ValorTotal      float64 `json:"valor_total"`
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

