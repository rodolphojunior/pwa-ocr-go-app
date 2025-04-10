package routes

import (
    "html/template"	
    "net/http"

	"github.com/gorilla/mux"
	"pwaocr/internal/handlers"
	"pwaocr/internal/auth"
    "pwaocr/internal/db"

)

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes() *mux.Router {
    r := mux.NewRouter()

    // Middleware global de CORS
    r.Use(EnableCORS)

    // Criar a instância do AuthHandler DENTRO da função
    authHandler := auth.AuthHandler{
        DB:        db.Conn, // ou nil, se for setado no main
        JWTSecret: []byte("sua-chave-secreta"),
    }

    // Carregar templates das páginas HTML uma vez
    tmplIndex := template.Must(template.ParseFiles("frontend/index.html"))
    tmplImg2Txt := template.Must(template.ParseFiles("frontend/img2txt.html"))
    //tmplPerfil := template.Must(template.ParseFiles("frontend/perfil.html"))

    // Rota pública para página inicial (index.html)
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := tmplIndex.Execute(w, nil); err != nil {
            http.Error(w, "Erro ao renderizar página inicial", http.StatusInternalServerError)
        }
    }).Methods("GET")

    // Rotas de autenticação (públicas)
    r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
    r.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")

    // Rota de upload de imagem (permanece pública por enquanto)
    r.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")

    // Subroteador para rotas protegidas por JWT
    secured := r.PathPrefix("").Subrouter()
    secured.Use(auth.AuthMiddleware)  // aplica AuthMiddleware apenas nestas rotas

    // Rota protegida para página de notas processadas (img2txt.html)
    secured.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := tmplIndex.Execute(w, nil); err != nil {
            http.Error(w, "Erro ao renderizar página img2txt", http.StatusInternalServerError)
        }
    }).Methods("GET")

    // Rota protegida para página de perfil.html
//    secured.HandleFunc("/perfil", func(w http.ResponseWriter, r *http.Request) {
//    w.Header().Set("Content-Type", "text/html; charset=utf-8")
//    if err := tmplPerfil.Execute(w, nil); err != nil {
//        http.Error(w, "Erro ao renderizar página perfil", http.StatusInternalServerError)
//    }
//}).Methods("GET")

    // ✅ Rota pública para carregar o HTML de perfil (proteção via JS)
    r.HandleFunc("/perfil", func(w http.ResponseWriter, r *http.Request) {
	    w.Header().Set("Content-Type", "text/html")
	    http.ServeFile(w, r, "./frontend/perfil.html")
    }).Methods("GET")

    r.HandleFunc("/img2txt", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        tmplImg2Txt.Execute(w, nil)
    }).Methods("GET")

    // Rotas protegidas de API (notas fiscais)
    secured.HandleFunc("/notas", handlers.ListarNotasHandler).Methods("GET")
    secured.HandleFunc("/notas", handlers.DeleteNotasHandler).Methods("DELETE")

    // Servir arquivos estáticos da pasta frontend/ (CSS, JS, imagens, manifest, etc.)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/")))
    


    return r
}

