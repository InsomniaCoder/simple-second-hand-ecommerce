package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
)

// Helper function to create a test repository
func newTestBboltRepository(t *testing.T) ItemRepository {
	t.Helper()
	// Create a unique temporary file for each test
	tmpFile, err := os.CreateTemp("", "bbolt-test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()

	repo, err := NewBboltRepository(dbPath)
	if err != nil {
		os.Remove(dbPath)
		t.Fatalf("failed to create test repository: %v", err)
	}
	t.Cleanup(func() {
		if closer, ok := repo.(interface{ Close() error }); ok {
			closer.Close()
		}
		os.Remove(dbPath)
	})
	return repo
}

func TestBboltRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		item    *models.Item
		setup   func(ItemRepository)
		wantErr bool
		errType error
	}{
		{
			name: "create new item",
			item: &models.Item{
				ID:          "test-1",
				Title:       "Test Item",
				Description: "Test Description",
				Price:       10.99,
				Status:      models.StatusAvailable,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "create item with minimal fields",
			item: &models.Item{
				ID:    "test-2",
				Title: "Minimal Item",
				Price: 5.00,
			},
			wantErr: false,
		},
		{
			name: "create duplicate item",
			item: &models.Item{
				ID:    "test-3",
				Title: "Duplicate",
				Price: 10.00,
			},
			setup: func(repo ItemRepository) {
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-3",
					Title: "Original",
					Price: 20.00,
				})
			},
			wantErr: true,
			errType: ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newTestBboltRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}

			err := repo.Create(context.Background(), tt.item)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("expected error type %v, got %v", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				// Verify item was stored
				stored, err := repo.GetByID(context.Background(), tt.item.ID)
				if err != nil {
					t.Errorf("failed to retrieve stored item: %v", err)
					return
				}
				if stored.ID != tt.item.ID {
					t.Errorf("stored item ID = %v, want %v", stored.ID, tt.item.ID)
				}
				if stored.Title != tt.item.Title {
					t.Errorf("stored item Title = %v, want %v", stored.Title, tt.item.Title)
				}
			}
		})
	}
}

func TestBboltRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		setup   func(ItemRepository)
		want    *models.Item
		wantErr bool
		errType error
	}{
		{
			name: "get existing item",
			id:   "test-1",
			setup: func(repo ItemRepository) {
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-1",
					Title: "Test Item",
					Price: 10.99,
				})
			},
			want: &models.Item{
				ID:    "test-1",
				Title: "Test Item",
				Price: 10.99,
			},
			wantErr: false,
		},
		{
			name:    "get non-existent item",
			id:      "non-existent",
			wantErr: true,
			errType: ErrNotFound,
		},
		{
			name: "get item from multiple items",
			id:   "test-2",
			setup: func(repo ItemRepository) {
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-1",
					Title: "Item 1",
					Price: 10.00,
				})
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-2",
					Title: "Item 2",
					Price: 20.00,
				})
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-3",
					Title: "Item 3",
					Price: 30.00,
				})
			},
			want: &models.Item{
				ID:    "test-2",
				Title: "Item 2",
				Price: 20.00,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newTestBboltRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}

			got, err := repo.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("expected error type %v, got %v", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if got.ID != tt.want.ID {
					t.Errorf("got ID = %v, want %v", got.ID, tt.want.ID)
				}
				if got.Title != tt.want.Title {
					t.Errorf("got Title = %v, want %v", got.Title, tt.want.Title)
				}
				if got.Price != tt.want.Price {
					t.Errorf("got Price = %v, want %v", got.Price, tt.want.Price)
				}
			}
		})
	}
}

func TestBboltRepository_List(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ItemRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:      "empty repository",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "single item",
			setup: func(repo ItemRepository) {
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-1",
					Title: "Item 1",
					Price: 10.00,
				})
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "multiple items",
			setup: func(repo ItemRepository) {
				for i := 1; i <= 5; i++ {
					_ = repo.Create(context.Background(), &models.Item{
						ID:    fmt.Sprintf("test-%d", i),
						Title: "Item",
						Price: float64(i * 10),
					})
				}
			},
			wantCount: 5,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newTestBboltRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}

			got, err := repo.List(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(got) != tt.wantCount {
				t.Errorf("got %d items, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestBboltRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		item    *models.Item
		setup   func(ItemRepository)
		wantErr bool
		errType error
	}{
		{
			name: "update existing item",
			item: &models.Item{
				ID:    "test-1",
				Title: "Updated Title",
				Price: 99.99,
			},
			setup: func(repo ItemRepository) {
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-1",
					Title: "Original Title",
					Price: 50.00,
				})
			},
			wantErr: false,
		},
		{
			name: "update non-existent item",
			item: &models.Item{
				ID:    "non-existent",
				Title: "Updated",
				Price: 10.00,
			},
			wantErr: true,
			errType: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newTestBboltRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}

			err := repo.Update(context.Background(), tt.item)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("expected error type %v, got %v", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				// Verify update was applied
				updated, err := repo.GetByID(context.Background(), tt.item.ID)
				if err != nil {
					t.Errorf("failed to retrieve updated item: %v", err)
					return
				}
				if updated.Title != tt.item.Title {
					t.Errorf("updated Title = %v, want %v", updated.Title, tt.item.Title)
				}
				if updated.Price != tt.item.Price {
					t.Errorf("updated Price = %v, want %v", updated.Price, tt.item.Price)
				}
			}
		})
	}
}

func TestBboltRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		setup   func(ItemRepository)
		wantErr bool
		errType error
	}{
		{
			name: "delete existing item",
			id:   "test-1",
			setup: func(repo ItemRepository) {
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-1",
					Title: "To Delete",
					Price: 10.00,
				})
			},
			wantErr: false,
		},
		{
			name:    "delete non-existent item",
			id:      "non-existent",
			wantErr: true,
			errType: ErrNotFound,
		},
		{
			name: "delete from multiple items",
			id:   "test-2",
			setup: func(repo ItemRepository) {
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-1",
					Title: "Item 1",
					Price: 10.00,
				})
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-2",
					Title: "Item 2",
					Price: 20.00,
				})
				_ = repo.Create(context.Background(), &models.Item{
					ID:    "test-3",
					Title: "Item 3",
					Price: 30.00,
				})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newTestBboltRepository(t)
			if tt.setup != nil {
				tt.setup(repo)
			}

			err := repo.Delete(context.Background(), tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
					return
				}
				if tt.errType != nil && !errors.Is(err, tt.errType) {
					t.Errorf("expected error type %v, got %v", tt.errType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}

				// Verify item was deleted
				_, err := repo.GetByID(context.Background(), tt.id)
				if !errors.Is(err, ErrNotFound) {
					t.Error("expected item to be deleted (ErrNotFound), but it still exists")
				}

				// If we had multiple items, verify others still exist
				if tt.name == "delete from multiple items" {
					items, _ := repo.List(context.Background())
					if len(items) != 2 {
						t.Errorf("expected 2 remaining items, got %d", len(items))
					}
				}
			}
		})
	}
}

// Concurrency Tests

func TestBboltRepository_ConcurrentCreate(t *testing.T) {
	repo := newTestBboltRepository(t)
	const numGoroutines = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			item := &models.Item{
				ID:    fmt.Sprintf("test-%d", id),
				Title: "Concurrent Item",
				Price: float64(id),
			}

			if err := repo.Create(context.Background(), item); err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("concurrent create error: %v", err)
	}

	// Verify all items were created
	items, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("failed to list items: %v", err)
	}
	if len(items) != numGoroutines {
		t.Errorf("expected %d items, got %d", numGoroutines, len(items))
	}
}

func TestBboltRepository_ConcurrentRead(t *testing.T) {
	repo := newTestBboltRepository(t)

	// Setup test items
	for i := 0; i < 10; i++ {
		_ = repo.Create(context.Background(), &models.Item{
			ID:    fmt.Sprintf("test-%d", i),
			Title: "Item",
			Price: float64(i),
		})
	}

	const numReaders = 50
	var wg sync.WaitGroup
	wg.Add(numReaders)

	errors := make(chan error, numReaders)

	for i := 0; i < numReaders; i++ {
		go func(readerID int) {
			defer wg.Done()

			id := fmt.Sprintf("test-%d", readerID%10)
			_, err := repo.GetByID(context.Background(), id)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("concurrent read error: %v", err)
	}
}

func TestBboltRepository_ConcurrentReadWrite(t *testing.T) {
	repo := newTestBboltRepository(t)

	// Setup initial items
	for i := 0; i < 10; i++ {
		_ = repo.Create(context.Background(), &models.Item{
			ID:    fmt.Sprintf("test-%d", i),
			Title: "Item",
			Price: float64(i),
		})
	}

	const numOperations = 100
	var wg sync.WaitGroup
	wg.Add(numOperations)

	for i := 0; i < numOperations; i++ {
		go func(opID int) {
			defer wg.Done()

			id := fmt.Sprintf("test-%d", opID%10)

			// Mix of operations
			switch opID % 3 {
			case 0: // Read
				_, _ = repo.GetByID(context.Background(), id)
			case 1: // Update
				_ = repo.Update(context.Background(), &models.Item{
					ID:    id,
					Title: "Updated",
					Price: float64(opID),
				})
			case 2: // List
				_, _ = repo.List(context.Background())
			}
		}(i)
	}

	wg.Wait()

	// Verify repository is still functional
	items, err := repo.List(context.Background())
	if err != nil {
		t.Errorf("repository corrupted after concurrent operations: %v", err)
	}
	if len(items) != 10 {
		t.Errorf("expected 10 items after concurrent operations, got %d", len(items))
	}
}

func TestBboltRepository_ConcurrentUpdate(t *testing.T) {
	repo := newTestBboltRepository(t)

	// Create single item
	_ = repo.Create(context.Background(), &models.Item{
		ID:    "test-1",
		Title: "Original",
		Price: 100.00,
	})

	const numUpdaters = 50
	var wg sync.WaitGroup
	wg.Add(numUpdaters)

	for i := 0; i < numUpdaters; i++ {
		go func(updaterID int) {
			defer wg.Done()

			_ = repo.Update(context.Background(), &models.Item{
				ID:    "test-1",
				Title: "Updated",
				Price: float64(updaterID),
			})
		}(i)
	}

	wg.Wait()

	// Verify item still exists and has a valid price
	item, err := repo.GetByID(context.Background(), "test-1")
	if err != nil {
		t.Errorf("item lost after concurrent updates: %v", err)
	}
	if item.Price < 0 || item.Price >= float64(numUpdaters) {
		t.Errorf("invalid price after concurrent updates: %v", item.Price)
	}
}

// File Persistence Tests

func TestBboltRepository_FilePersistence(t *testing.T) {
	// Create temporary file for database
	tmpFile, err := os.CreateTemp("", "bbolt-test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(dbPath)

	// Create repository and add item
	repo, err := NewBboltRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	testItem := &models.Item{
		ID:    "persist-1",
		Title: "Persistent Item",
		Price: 99.99,
	}

	err = repo.Create(context.Background(), testItem)
	if err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	// Close repository
	if closer, ok := repo.(interface{ Close() error }); ok {
		closer.Close()
	}

	// Reopen repository
	repo2, err := NewBboltRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to reopen repository: %v", err)
	}
	defer func() {
		if closer, ok := repo2.(interface{ Close() error }); ok {
			closer.Close()
		}
	}()

	// Verify item persisted
	retrieved, err := repo2.GetByID(context.Background(), "persist-1")
	if err != nil {
		t.Fatalf("failed to retrieve persisted item: %v", err)
	}

	if retrieved.ID != testItem.ID {
		t.Errorf("persisted item ID = %v, want %v", retrieved.ID, testItem.ID)
	}
	if retrieved.Title != testItem.Title {
		t.Errorf("persisted item Title = %v, want %v", retrieved.Title, testItem.Title)
	}
	if retrieved.Price != testItem.Price {
		t.Errorf("persisted item Price = %v, want %v", retrieved.Price, testItem.Price)
	}
}

func TestBboltRepository_Close(t *testing.T) {
	repo, err := NewBboltRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	// Verify repository implements Close method
	closer, ok := repo.(interface{ Close() error })
	if !ok {
		t.Fatal("repository does not implement Close() error")
	}

	// Close repository
	if err := closer.Close(); err != nil {
		t.Errorf("failed to close repository: %v", err)
	}
}
