package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
)

var (
	defaultBucketName = []byte("items")
)

// bboltRepository implements ItemRepository with bbolt storage
type bboltRepository struct {
	db         *bbolt.DB
	bucketName []byte
}

// NewBboltRepository creates a new bbolt-based repository
// path can be a file path for persistent storage or ":memory:" for in-memory mode
func NewBboltRepository(path string) (ItemRepository, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bbolt database: %w", err)
	}

	repo := &bboltRepository{
		db:         db,
		bucketName: defaultBucketName,
	}

	// Create bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(repo.bucketName)
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return repo, nil
}

// Close closes the bbolt database
func (r *bboltRepository) Close() error {
	return r.db.Close()
}

// Create adds a new item to the repository
func (r *bboltRepository) Create(ctx context.Context, item *models.Item) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(r.bucketName)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// Check if item already exists
		existing := b.Get([]byte(item.ID))
		if existing != nil {
			return ErrAlreadyExists
		}

		// Marshal item to JSON
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to marshal item: %w", err)
		}

		// Store item
		if err := b.Put([]byte(item.ID), data); err != nil {
			return fmt.Errorf("failed to store item: %w", err)
		}

		return nil
	})
}

// GetByID retrieves an item by its ID
func (r *bboltRepository) GetByID(ctx context.Context, id string) (*models.Item, error) {
	var item *models.Item

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(r.bucketName)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		data := b.Get([]byte(id))
		if data == nil {
			return ErrNotFound
		}

		item = &models.Item{}
		if err := json.Unmarshal(data, item); err != nil {
			return fmt.Errorf("failed to unmarshal item: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return item, nil
}

// List returns all items
func (r *bboltRepository) List(ctx context.Context) ([]*models.Item, error) {
	var items []*models.Item

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(r.bucketName)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		return b.ForEach(func(k, v []byte) error {
			var item models.Item
			if err := json.Unmarshal(v, &item); err != nil {
				return fmt.Errorf("failed to unmarshal item: %w", err)
			}
			items = append(items, &item)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// Update modifies an existing item
func (r *bboltRepository) Update(ctx context.Context, item *models.Item) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(r.bucketName)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// Check if item exists
		existing := b.Get([]byte(item.ID))
		if existing == nil {
			return ErrNotFound
		}

		// Marshal item to JSON
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to marshal item: %w", err)
		}

		// Update item
		if err := b.Put([]byte(item.ID), data); err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}

		return nil
	})
}

// Delete removes an item from the repository
func (r *bboltRepository) Delete(ctx context.Context, id string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(r.bucketName)
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// Check if item exists
		existing := b.Get([]byte(id))
		if existing == nil {
			return ErrNotFound
		}

		// Delete item
		if err := b.Delete([]byte(id)); err != nil {
			return fmt.Errorf("failed to delete item: %w", err)
		}

		return nil
	})
}
