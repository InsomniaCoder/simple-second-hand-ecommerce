package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/repository"
	"github.com/google/uuid"
)

type itemService struct {
	repo repository.ItemRepository
}

// NewItemService creates a new item service
func NewItemService(repo repository.ItemRepository) ItemService {
	return &itemService{repo: repo}
}

// CreateItem creates a new item with validation
func (s *itemService) CreateItem(ctx context.Context, item *models.Item) error {
	if item == nil {
		return fmt.Errorf("%w: item cannot be nil", ErrInvalidInput)
	}

	// Validate item
	if err := item.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Set defaults
	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	// Ensure default status if not set
	if item.Status == "" {
		item.Status = models.StatusAvailable
	}

	// Create in repository
	if err := s.repo.Create(ctx, item); err != nil {
		return fmt.Errorf("%w: %v", ErrCreateFailed, err)
	}

	return nil
}

// GetItem retrieves an item by ID
func (s *itemService) GetItem(ctx context.Context, id string) (*models.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return item, nil
}

// ListItems returns all items
func (s *itemService) ListItems(ctx context.Context) ([]*models.Item, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	return items, nil
}

// UpdateItem updates an existing item
func (s *itemService) UpdateItem(ctx context.Context, id string, updates *models.Item) error {
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidInput)
	}
	if updates == nil {
		return fmt.Errorf("%w: updates cannot be nil", ErrInvalidInput)
	}

	// Get existing item
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrItemNotFound
		}
		return fmt.Errorf("failed to get item: %w", err)
	}

	// Apply updates
	if updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Description != "" {
		existing.Description = updates.Description
	}
	if updates.Price > 0 {
		existing.Price = updates.Price
	}
	if updates.Image != "" {
		existing.Image = updates.Image
	}
	if updates.Status != "" {
		existing.Status = updates.Status
	}

	// Validate merged item
	if err := existing.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidInput, err)
	}

	// Update timestamp
	existing.UpdatedAt = time.Now()

	// Update in repository
	if err := s.repo.Update(ctx, existing); err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateFailed, err)
	}

	return nil
}

// DeleteItem removes an item
func (s *itemService) DeleteItem(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrInvalidInput)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrItemNotFound
		}
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}

	return nil
}
