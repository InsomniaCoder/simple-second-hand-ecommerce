package handlers

import (
	"net/http"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/service"
)

// SetupRoutes configures all HTTP routes with middleware
func SetupRoutes(svc service.ItemService) http.Handler {
	mux := http.NewServeMux()
	handlers := NewItemHandlers(svc)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Item endpoints
	mux.HandleFunc("POST /items", handlers.CreateItem)
	mux.HandleFunc("GET /items", handlers.ListItems)
	mux.HandleFunc("GET /items/{id}", handlers.GetItem)
	mux.HandleFunc("PUT /items/{id}", handlers.UpdateItem)
	mux.HandleFunc("DELETE /items/{id}", handlers.DeleteItem)

	// Apply middleware
	return applyMiddleware(mux)
}

// applyMiddleware wraps the handler with middleware
func applyMiddleware(handler http.Handler) http.Handler {
	// Apply middleware in reverse order (last middleware wraps first)
	handler = ContentTypeJSON(handler)
	handler = CORS(handler)
	handler = Logger(handler)
	handler = Recovery(handler)
	return handler
}
