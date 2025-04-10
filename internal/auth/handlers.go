// internal/auth/handlers.go
package auth

import (
	"encoding/json"
	"net/http"
//	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	JWTSecret []byte
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao gerar hash", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashed)
	if err := h.DB.Create(&user).Error; err != nil {
		http.Error(w, "Erro ao registrar usuário", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}
	var user User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString(h.JWTSecret)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}

func (h *AuthHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey)
	if userID == nil {
		http.Error(w, "Usuário não autenticado", http.StatusUnauthorized)
		return
	}

	var user User
	if err := h.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	// Retorna apenas os dados não sensíveis
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

