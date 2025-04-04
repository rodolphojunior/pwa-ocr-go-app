// internal/http/routes.go
package http

import (
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("/upload", UploadHandler)
}

