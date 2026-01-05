package service

import (
	"context"
	"errors"
	"testing"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/repository"
)

// mockRepository implements ItemRepository for testing
type mockRepository struct {
	CreateFunc  func(ctx context.Context, item *models.Item) error
	GetByIDFunc func(ctx context.Context, id string) (*models.Item, error)
	ListFunc    func(ctx context.Context) ([]*models.Item, error)
	UpdateFunc  func(ctx context.Context, item *models.Item) error
	DeleteFunc  func(ctx context.Context, id string) error
}

func (m *mockRepository) Create(ctx context.Context, item *models.Item) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, item)
	}
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*models.Item, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, repository.ErrNotFound
}

func (m *mockRepository) List(ctx context.Context) ([]*models.Item, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return []*models.Item{}, nil
}

func (m *mockRepository) Update(ctx context.Context, item *models.Item) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, item)
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestItemService_CreateItem(t *testing.T) {
	tests := []struct {
		name    string
		item    *models.Item
		setup   func(*mockRepository)
		wantErr bool
		errType error
		verify  func(*testing.T, *models.Item)
	}{
		{
			name: "create valid item",
			item: &models.Item{
				Title:       "Test Item",
				Description: "Test Description",
				Price:       10.99,
				Status:      models.StatusAvailable,
			},
			setup: func(m *mockRepository) {
				m.CreateFunc = func(ctx context.Context, item *models.Item) error {
					return nil
				}
			},
			wantErr: false,
			verify: func(t *testing.T, item *models.Item) {
				if item.ID == "" {
					t.Error("expected ID to be generated")
				}
				if item.CreatedAt.IsZero() {
					t.Error("expected CreatedAt to be set")
				}
				if item.UpdatedAt.IsZero() {
					t.Error("expected UpdatedAt to be set")
				}
			},
		},
		{
			name: "create item with explicit status",
			item: &models.Item{
				Title:  "Item with status",
				Price:  5.00,
				Status: models.StatusBooked,
			},
			setup: func(m *mockRepository) {
				m.CreateFunc = func(ctx context.Context, item *models.Item) error {
					return nil
				}
			},
			wantErr: false,
			verify: func(t *testing.T, item *models.Item) {
				if item.Status != models.StatusBooked {
					t.Errorf("expected status %v, got %v", models.StatusBooked, item.Status)
				}
			},
		},
		{
			name:    "create nil item",
			item:    nil,
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name: "create item with empty title",
			item: &models.Item{
				Title: "",
				Price: 10.00,
			},
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name: "create item with invalid price",
			item: &models.Item{
				Title: "Invalid Price",
				Price: -5.00,
			},
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name: "create item with title too long",
			item: &models.Item{
				Title: string(make([]byte, 201)),
				Price: 10.00,
			},
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name: "repository error",
			item: &models.Item{
				Title:  "Test",
				Price:  10.00,
				Status: models.StatusAvailable,
			},
			setup: func(m *mockRepository) {
				m.CreateFunc = func(ctx context.Context, item *models.Item) error {
					return errors.New("repository error")
				}
			},
			wantErr: true,
			errType: ErrCreateFailed,
		},
		{
			name: "repository already exists error",
			item: &models.Item{
				Title:  "Duplicate",
				Price:  10.00,
				Status: models.StatusAvailable,
			},
			setup: func(m *mockRepository) {
				m.CreateFunc = func(ctx context.Context, item *models.Item) error {
					return repository.ErrAlreadyExists
				}
			},
			wantErr: true,
			errType: ErrCreateFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			svc := NewItemService(mockRepo)
			err := svc.CreateItem(context.Background(), tt.item)

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
				if tt.verify != nil {
					tt.verify(t, tt.item)
				}
			}
		})
	}
}

func TestItemService_GetItem(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		setup   func(*mockRepository)
		want    *models.Item
		wantErr bool
		errType error
	}{
		{
			name: "get existing item",
			id:   "test-1",
			setup: func(m *mockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return &models.Item{
						ID:    "test-1",
						Title: "Test Item",
						Price: 10.99,
					}, nil
				}
			},
			want: &models.Item{
				ID:    "test-1",
				Title: "Test Item",
				Price: 10.99,
			},
			wantErr: false,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name: "item not found",
			id:   "non-existent",
			setup: func(m *mockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return nil, repository.ErrNotFound
				}
			},
			wantErr: true,
			errType: ErrItemNotFound,
		},
		{
			name: "repository error",
			id:   "test-1",
			setup: func(m *mockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			svc := NewItemService(mockRepo)
			got, err := svc.GetItem(context.Background(), tt.id)

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
			}
		})
	}
}

func TestItemService_ListItems(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name: "list empty",
			setup: func(m *mockRepository) {
				m.ListFunc = func(ctx context.Context) ([]*models.Item, error) {
					return []*models.Item{}, nil
				}
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "list multiple items",
			setup: func(m *mockRepository) {
				m.ListFunc = func(ctx context.Context) ([]*models.Item, error) {
					return []*models.Item{
						{ID: "1", Title: "Item 1", Price: 10.00},
						{ID: "2", Title: "Item 2", Price: 20.00},
						{ID: "3", Title: "Item 3", Price: 30.00},
					}, nil
				}
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "repository error",
			setup: func(m *mockRepository) {
				m.ListFunc = func(ctx context.Context) ([]*models.Item, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			svc := NewItemService(mockRepo)
			got, err := svc.ListItems(context.Background())

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

func TestItemService_UpdateItem(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		updates *models.Item
		setup   func(*mockRepository)
		wantErr bool
		errType error
		verify  func(*testing.T, *mockRepository)
	}{
		{
			name: "update title",
			id:   "test-1",
			updates: &models.Item{
				Title: "Updated Title",
			},
			setup: func(m *mockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return &models.Item{
						ID:     "test-1",
						Title:  "Original Title",
						Price:  10.00,
						Status: models.StatusAvailable,
					}, nil
				}
				m.UpdateFunc = func(ctx context.Context, item *models.Item) error {
					return nil
				}
			},
			wantErr: false,
			verify: func(t *testing.T, m *mockRepository) {
				// Verify update was called
			},
		},
		{
			name: "update price",
			id:   "test-1",
			updates: &models.Item{
				Price: 99.99,
			},
			setup: func(m *mockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return &models.Item{
						ID:     "test-1",
						Title:  "Item",
						Price:  10.00,
						Status: models.StatusAvailable,
					}, nil
				}
				m.UpdateFunc = func(ctx context.Context, item *models.Item) error {
					if item.Price != 99.99 {
						return errors.New("price not updated")
					}
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "update multiple fields",
			id:   "test-1",
			updates: &models.Item{
				Title:  "New Title",
				Price:  50.00,
				Status: models.StatusBooked,
			},
			setup: func(m *mockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return &models.Item{
						ID:     "test-1",
						Title:  "Old Title",
						Price:  10.00,
						Status: models.StatusAvailable,
					}, nil
				}
				m.UpdateFunc = func(ctx context.Context, item *models.Item) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:    "empty id",
			id:      "",
			updates: &models.Item{Title: "Test"},
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name:    "nil updates",
			id:      "test-1",
			updates: nil,
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name: "item not found",
			id:   "non-existent",
			updates: &models.Item{
				Title: "Update",
			},
			setup: func(m *mockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return nil, repository.ErrNotFound
				}
			},
			wantErr: true,
			errType: ErrItemNotFound,
		},
		{
			name: "update makes item invalid - invalid status",
			id:   "test-1",
			updates: &models.Item{
				Status: "invalid-status",
			},
			setup: func(m *mockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return &models.Item{
						ID:     "test-1",
						Title:  "Item",
						Price:  10.00,
						Status: models.StatusAvailable,
					}, nil
				}
			},
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name: "repository update error",
			id:   "test-1",
			updates: &models.Item{
				Title: "Updated",
			},
			setup: func(m *mockRepository) {
				m.GetByIDFunc = func(ctx context.Context, id string) (*models.Item, error) {
					return &models.Item{
						ID:     "test-1",
						Title:  "Item",
						Price:  10.00,
						Status: models.StatusAvailable,
					}, nil
				}
				m.UpdateFunc = func(ctx context.Context, item *models.Item) error {
					return errors.New("database error")
				}
			},
			wantErr: true,
			errType: ErrUpdateFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			svc := NewItemService(mockRepo)
			err := svc.UpdateItem(context.Background(), tt.id, tt.updates)

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
				if tt.verify != nil {
					tt.verify(t, mockRepo)
				}
			}
		})
	}
}

func TestItemService_DeleteItem(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		setup   func(*mockRepository)
		wantErr bool
		errType error
	}{
		{
			name: "delete existing item",
			id:   "test-1",
			setup: func(m *mockRepository) {
				m.DeleteFunc = func(ctx context.Context, id string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name: "item not found",
			id:   "non-existent",
			setup: func(m *mockRepository) {
				m.DeleteFunc = func(ctx context.Context, id string) error {
					return repository.ErrNotFound
				}
			},
			wantErr: true,
			errType: ErrItemNotFound,
		},
		{
			name: "repository error",
			id:   "test-1",
			setup: func(m *mockRepository) {
				m.DeleteFunc = func(ctx context.Context, id string) error {
					return errors.New("database error")
				}
			},
			wantErr: true,
			errType: ErrDeleteFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			if tt.setup != nil {
				tt.setup(mockRepo)
			}

			svc := NewItemService(mockRepo)
			err := svc.DeleteItem(context.Background(), tt.id)

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
			}
		})
	}
}
