package service

import (
	"context"
	"errors"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
)

// Service errors
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrItemNotFound = errors.New("item not found")
	ErrUpdateFailed = errors.New("failed to update item")
	ErrDeleteFailed = errors.New("failed to delete item")
	ErrCreateFailed = errors.New("failed to create item")
)

// ItemService defines business logic operations for items
type ItemService interface {
	CreateItem(ctx context.Context, item *models.Item) error
	GetItem(ctx context.Context, id string) (*models.Item, error)
	ListItems(ctx context.Context) ([]*models.Item, error)
	UpdateItem(ctx context.Context, id string, updates *models.Item) error
	DeleteItem(ctx context.Context, id string) error
}
