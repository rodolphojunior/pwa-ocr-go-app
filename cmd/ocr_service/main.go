package main

import (
	"log"
    "net"
	"net/http"
    "github.com/joho/godotenv"
//	"pwaocr/internal/auth"
    "pwaocr/internal/routes"
	"pwaocr/internal/db"
)

func EnableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "localhost"
}

func main() {

    // Carrega as variáveis do .env
    if err := godotenv.Load(); err != nil {
        log.Println("Aviso: .env não carregado (pode já estar definido no sistema)")
    }

	db.Connect() // conecta ao banco de dados

    // Atualiza o AuthHandler com a conexão real
//    authHandler := auth.AuthHandler{
//        DB: db.Conn,
//        JWTSecret: []byte("sua-chave-secreta"),
//    }

    db.RunMigrations() // roda as migrações todas (veja db.go)

	router := routes.SetupRoutes()

    ip := getLocalIP()
    log.Printf("Servidor disponível em:\n- Local: http://localhost:8080\n- Celular: http://%s:8080\n", ip)

    log.Fatal(http.ListenAndServe("0.0.0.0:8080", EnableCORS(router)))

}

