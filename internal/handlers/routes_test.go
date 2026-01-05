package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
)

func TestSetupRoutes(t *testing.T) {
	mockSvc := &mockService{
		ListItemsFunc: func(ctx context.Context) ([]*models.Item, error) {
			return []*models.Item{}, nil
		},
	}

	router := SetupRoutes(mockSvc)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "health endpoint",
			method:     http.MethodGet,
			path:       "/health",
			wantStatus: http.StatusOK,
		},
		{
			name:       "list items endpoint",
			method:     http.MethodGet,
			path:       "/items",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestApplyMiddleware(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := applyMiddleware(mux)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify middleware was applied (check for CORS and Content-Type headers)
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS middleware to be applied")
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("expected ContentTypeJSON middleware to be applied")
	}
}
