package repository

import (
	"context"
	"errors"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
)

// Common repository errors
var (
	ErrNotFound      = errors.New("item not found")
	ErrAlreadyExists = errors.New("item already exists")
)

// ItemRepository defines the interface for item data access
type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
	GetByID(ctx context.Context, id string) (*models.Item, error)
	List(ctx context.Context) ([]*models.Item, error)
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, id string) error
}
