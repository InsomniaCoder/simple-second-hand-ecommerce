package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/service"
)

// ItemHandlers handles HTTP requests for items
type ItemHandlers struct {
	service service.ItemService
}

// NewItemHandlers creates new item handlers
func NewItemHandlers(service service.ItemService) *ItemHandlers {
	return &ItemHandlers{service: service}
}

// CreateItemRequest represents the create item request body
type CreateItemRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	Image       string            `json:"image"`
	Status      models.ItemStatus `json:"status,omitempty"`
}

// UpdateItemRequest represents the update item request body
type UpdateItemRequest struct {
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	Price       float64           `json:"price,omitempty"`
	Image       string            `json:"image,omitempty"`
	Status      models.ItemStatus `json:"status,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// respondJSON writes JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// respondError writes error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}

// CreateItem handles POST /items
func (h *ItemHandlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item := &models.Item{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Image:       req.Image,
		Status:      req.Status,
	}

	if err := h.service.CreateItem(r.Context(), item); err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, item)
}

// ListItems handles GET /items
func (h *ItemHandlers) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.ListItems(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, items)
}

// GetItem handles GET /items/{id}
func (h *ItemHandlers) GetItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}

	item, err := h.service.GetItem(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			respondError(w, http.StatusNotFound, "item not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, item)
}

// UpdateItem handles PUT /items/{id}
func (h *ItemHandlers) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}

	var req UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updates := &models.Item{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Image:       req.Image,
		Status:      req.Status,
	}

	if err := h.service.UpdateItem(r.Context(), id, updates); err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			respondError(w, http.StatusNotFound, "item not found")
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	item, _ := h.service.GetItem(r.Context(), id)
	respondJSON(w, http.StatusOK, item)
}

// DeleteItem handles DELETE /items/{id}
func (h *ItemHandlers) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.service.DeleteItem(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrItemNotFound) {
			respondError(w, http.StatusNotFound, "item not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "item deleted"})
}
