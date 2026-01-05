package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/service"
)

// mockService implements ItemService for testing
type mockService struct {
	CreateItemFunc func(ctx context.Context, item *models.Item) error
	GetItemFunc    func(ctx context.Context, id string) (*models.Item, error)
	ListItemsFunc  func(ctx context.Context) ([]*models.Item, error)
	UpdateItemFunc func(ctx context.Context, id string, updates *models.Item) error
	DeleteItemFunc func(ctx context.Context, id string) error
}

func (m *mockService) CreateItem(ctx context.Context, item *models.Item) error {
	if m.CreateItemFunc != nil {
		return m.CreateItemFunc(ctx, item)
	}
	return nil
}

func (m *mockService) GetItem(ctx context.Context, id string) (*models.Item, error) {
	if m.GetItemFunc != nil {
		return m.GetItemFunc(ctx, id)
	}
	return nil, service.ErrItemNotFound
}

func (m *mockService) ListItems(ctx context.Context) ([]*models.Item, error) {
	if m.ListItemsFunc != nil {
		return m.ListItemsFunc(ctx)
	}
	return []*models.Item{}, nil
}

func (m *mockService) UpdateItem(ctx context.Context, id string, updates *models.Item) error {
	if m.UpdateItemFunc != nil {
		return m.UpdateItemFunc(ctx, id, updates)
	}
	return nil
}

func (m *mockService) DeleteItem(ctx context.Context, id string) error {
	if m.DeleteItemFunc != nil {
		return m.DeleteItemFunc(ctx, id)
	}
	return nil
}

func TestItemHandlers_CreateItem(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setup          func(*mockService)
		wantStatus     int
		wantErrMessage string
	}{
		{
			name: "create valid item",
			body: CreateItemRequest{
				Title:       "Test Item",
				Description: "Test Description",
				Price:       10.99,
				Image:       "https://example.com/image.jpg",
				Status:      models.StatusAvailable,
			},
			setup: func(m *mockService) {
				m.CreateItemFunc = func(ctx context.Context, item *models.Item) error {
					item.ID = "test-123"
					return nil
				}
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid JSON body",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service validation error",
			body: CreateItemRequest{
				Title: "",
				Price: 10.00,
			},
			setup: func(m *mockService) {
				m.CreateItemFunc = func(ctx context.Context, item *models.Item) error {
					return service.ErrInvalidInput
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service internal error",
			body: CreateItemRequest{
				Title: "Test",
				Price: 10.00,
			},
			setup: func(m *mockService) {
				m.CreateItemFunc = func(ctx context.Context, item *models.Item) error {
					return errors.New("database error")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockService{}
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			handlers := NewItemHandlers(mockSvc)

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handlers.CreateItem(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusCreated {
				var item models.Item
				if err := json.NewDecoder(w.Body).Decode(&item); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
			}
		})
	}
}

func TestItemHandlers_ListItems(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*mockService)
		wantStatus int
		wantCount  int
	}{
		{
			name: "list empty",
			setup: func(m *mockService) {
				m.ListItemsFunc = func(ctx context.Context) ([]*models.Item, error) {
					return []*models.Item{}, nil
				}
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "list multiple items",
			setup: func(m *mockService) {
				m.ListItemsFunc = func(ctx context.Context) ([]*models.Item, error) {
					return []*models.Item{
						{ID: "1", Title: "Item 1", Price: 10.00},
						{ID: "2", Title: "Item 2", Price: 20.00},
					}, nil
				}
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "service error",
			setup: func(m *mockService) {
				m.ListItemsFunc = func(ctx context.Context) ([]*models.Item, error) {
					return nil, errors.New("database error")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockService{}
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			handlers := NewItemHandlers(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/items", nil)
			w := httptest.NewRecorder()

			handlers.ListItems(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var items []*models.Item
				if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if len(items) != tt.wantCount {
					t.Errorf("got %d items, want %d", len(items), tt.wantCount)
				}
			}
		})
	}
}

func TestItemHandlers_GetItem(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setup      func(*mockService)
		wantStatus int
	}{
		{
			name: "get existing item",
			id:   "test-1",
			setup: func(m *mockService) {
				m.GetItemFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return &models.Item{
						ID:    id,
						Title: "Test Item",
						Price: 10.99,
					}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "item not found",
			id:   "non-existent",
			setup: func(m *mockService) {
				m.GetItemFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return nil, service.ErrItemNotFound
				}
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "service error",
			id:   "test-1",
			setup: func(m *mockService) {
				m.GetItemFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return nil, errors.New("database error")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockService{}
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			handlers := NewItemHandlers(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/items/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()

			handlers.GetItem(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var item models.Item
				if err := json.NewDecoder(w.Body).Decode(&item); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if item.ID != tt.id {
					t.Errorf("got ID %s, want %s", item.ID, tt.id)
				}
			}
		})
	}
}

func TestItemHandlers_UpdateItem(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		body       interface{}
		setup      func(*mockService)
		wantStatus int
	}{
		{
			name: "update item successfully",
			id:   "test-1",
			body: UpdateItemRequest{
				Title: "Updated Title",
				Price: 99.99,
			},
			setup: func(m *mockService) {
				m.UpdateItemFunc = func(ctx context.Context, id string, updates *models.Item) error {
					return nil
				}
				m.GetItemFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return &models.Item{
						ID:    id,
						Title: "Updated Title",
						Price: 99.99,
					}, nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid JSON body",
			id:         "test-1",
			body:       "invalid json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "item not found",
			id:   "non-existent",
			body: UpdateItemRequest{
				Title: "Updated",
			},
			setup: func(m *mockService) {
				m.UpdateItemFunc = func(ctx context.Context, id string, updates *models.Item) error {
					return service.ErrItemNotFound
				}
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "validation error",
			id:   "test-1",
			body: UpdateItemRequest{
				Price: -10.00,
			},
			setup: func(m *mockService) {
				m.UpdateItemFunc = func(ctx context.Context, id string, updates *models.Item) error {
					return service.ErrInvalidInput
				}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			id:   "test-1",
			body: UpdateItemRequest{
				Title: "Updated",
			},
			setup: func(m *mockService) {
				m.UpdateItemFunc = func(ctx context.Context, id string, updates *models.Item) error {
					return errors.New("database error")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockService{}
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			handlers := NewItemHandlers(mockSvc)

			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPut, "/items/"+tt.id, bytes.NewReader(bodyBytes))
			req.SetPathValue("id", tt.id)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handlers.UpdateItem(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestItemHandlers_DeleteItem(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		setup      func(*mockService)
		wantStatus int
	}{
		{
			name: "delete successfully",
			id:   "test-1",
			setup: func(m *mockService) {
				m.DeleteItemFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "item not found",
			id:   "non-existent",
			setup: func(m *mockService) {
				m.DeleteItemFunc = func(ctx context.Context, id string) error {
					return service.ErrItemNotFound
				}
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "service error",
			id:   "test-1",
			setup: func(m *mockService) {
				m.DeleteItemFunc = func(ctx context.Context, id string) error {
					return errors.New("database error")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockService{}
			if tt.setup != nil {
				tt.setup(mockSvc)
			}

			handlers := NewItemHandlers(mockSvc)

			req := httptest.NewRequest(http.MethodDelete, "/items/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()

			handlers.DeleteItem(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var response map[string]string
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response["message"] != "item deleted" {
					t.Errorf("got message %q, want %q", response["message"], "item deleted")
				}
			}
		})
	}
}

func TestRespondJSON(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		data       interface{}
		wantStatus int
	}{
		{
			name:       "respond with data",
			status:     http.StatusOK,
			data:       map[string]string{"message": "success"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "respond with nil data",
			status:     http.StatusNoContent,
			data:       nil,
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			respondJSON(w, tt.status, tt.data)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()
	respondError(w, http.StatusBadRequest, "test error")

	if w.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", w.Code, http.StatusBadRequest)
	}

	var errResp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error != "test error" {
		t.Errorf("got error %q, want %q", errResp.Error, "test error")
	}
}
