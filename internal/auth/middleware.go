// internal/auth/middleware.go
package auth

import (
	//"context"
	"net/http"
	"strings"
    "log"
    "fmt"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey = contextKey("userID")

//func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		authHeader := r.Header.Get("Authorization")
//		if authHeader == "" {
//			http.Error(w, "Token ausente", http.StatusUnauthorized)
//			return
//		}
//		parts := strings.Split(authHeader, " ")
//		if len(parts) != 2 || parts[0] != "Bearer" {
//			http.Error(w, "Formato do token inválido", http.StatusUnauthorized)
//			return
//		}
//		tokenStr := parts[1]
//		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
//			return h.JWTSecret, nil
//		})
//		if err != nil || !token.Valid {
//			http.Error(w, "Token inválido", http.StatusUnauthorized)
//			return
//		}
//		claims, ok := token.Claims.(jwt.MapClaims)
//		if !ok {
//			http.Error(w, "Token inválido", http.StatusUnauthorized)
//			return
//		}
//		userID := claims["user_id"]
//		ctx := context.WithValue(r.Context(), UserIDKey, userID)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	})
//}
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")

        log.Println("🛡 Recebido Authorization Header:", authHeader)

        if authHeader == "" {
            http.Error(w, `{"error": "Token ausente"}`, http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, `{"error": "Token malformatado"}`, http.StatusUnauthorized)
            return
        }

        tokenString := parts[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("método de assinatura inesperado")
            }
            return []byte("sua-chave-secreta"), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, `{"error": "Token inválido"}`, http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

