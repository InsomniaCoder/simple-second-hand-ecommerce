package repository

import (
	"context"
	"sync"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
)

// memoryRepository implements ItemRepository with in-memory storage
type memoryRepository struct {
	mu    sync.RWMutex
	items map[string]*models.Item
}

// NewMemoryRepository creates a new in-memory repository
func NewMemoryRepository() ItemRepository {
	return &memoryRepository{
		items: make(map[string]*models.Item),
	}
}

// Create adds a new item to the repository
func (r *memoryRepository) Create(ctx context.Context, item *models.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; exists {
		return ErrAlreadyExists
	}

	r.items[item.ID] = item
	return nil
}

// GetByID retrieves an item by its ID
func (r *memoryRepository) GetByID(ctx context.Context, id string) (*models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}

	return item, nil
}

// List returns all items
func (r *memoryRepository) List(ctx context.Context) ([]*models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*models.Item, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}

	return items, nil
}

// Update modifies an existing item
func (r *memoryRepository) Update(ctx context.Context, item *models.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[item.ID]; !exists {
		return ErrNotFound
	}

	r.items[item.ID] = item
	return nil
}

// Delete removes an item from the repository
func (r *memoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return ErrNotFound
	}

	delete(r.items, id)
	return nil
}
