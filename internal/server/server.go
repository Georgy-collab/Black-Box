package server

import (
	"net/http"
	"os"

	"github.com/Georgy-collab/Black-Box/internal/handlers"
	"github.com/Georgy-collab/Black-Box/internal/store"
)

func New() *http.ServeMux {
	userStore := store.NewUserStore()
	h := handlers.NewUserHandler(userStore)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("GET /users", h.ListUsers)
	mux.HandleFunc("GET /users/{user_id}", h.GetUser)
	mux.HandleFunc("PUT /users/{user_id}", h.UpdateUser)
	mux.HandleFunc("DELETE /users/{user_id}", h.DeleteUser)

	mux.HandleFunc("GET /openapi.json", serveFile("openapi.json", "application/json"))
	mux.HandleFunc("GET /openapi.yaml", serveFile("openapi.yaml", "application/yaml"))
	mux.HandleFunc("GET /docs", serveDocs)

	return mux
}

func serveFile(path, contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(path)
		if err != nil {
			http.Error(w, "OpenAPI file not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", contentType)
		_, _ = w.Write(data)
	}
}

func serveDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>User Service API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: '/openapi.json',
      dom_id: '#swagger-ui'
    });
  </script>
</body>
</html>`))
}
