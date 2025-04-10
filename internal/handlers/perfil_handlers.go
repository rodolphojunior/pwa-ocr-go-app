package handlers

import (
	"html/template"
	"net/http"
)

func PerfilHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("frontend/perfil.html")
	if err != nil {
		http.Error(w, "Erro ao carregar template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

