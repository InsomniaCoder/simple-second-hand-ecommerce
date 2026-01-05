package loader

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/repository"
)

func TestCSVLoader_LoadFromFile(t *testing.T) {
	tests := []struct {
		name          string
		csvContent    string
		expectedCount int
		wantErr       bool
		errContains   string
	}{
		{
			name: "valid CSV with 3 items",
			csvContent: `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,Vintage Leather Jacket,Classic brown leather jacket,150.00,https://example.com/jacket.jpg,available
550e8400-e29b-41d4-a716-446655440002,MacBook Pro 2019,15-inch MacBook Pro,800.00,https://example.com/macbook.jpg,available
550e8400-e29b-41d4-a716-446655440003,Ikea Standing Desk,Electric standing desk,250.00,https://example.com/desk.jpg,available`,
			expectedCount: 3,
			wantErr:       false,
		},
		{
			name: "valid CSV with 1 item",
			csvContent: `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,Test Item,Description,100.00,https://example.com/test.jpg,available`,
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name:          "empty CSV (header only)",
			csvContent:    `id,title,description,price,image,status`,
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "invalid header",
			csvContent: `wrong,header,format
550e8400-e29b-41d4-a716-446655440001,Test Item,Description,100.00,image.jpg,available`,
			expectedCount: 0,
			wantErr:       true,
			errContains:   "invalid CSV header",
		},
		{
			name: "invalid price format",
			csvContent: `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,Test Item,Description,invalid_price,image.jpg,available`,
			expectedCount: 0,
			wantErr:       true,
			errContains:   "invalid price value",
		},
		{
			name: "negative price",
			csvContent: `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,Test Item,Description,-10.00,image.jpg,available`,
			expectedCount: 0,
			wantErr:       true,
			errContains:   "price must be greater than 0",
		},
		{
			name: "invalid status",
			csvContent: `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,Test Item,Description,100.00,image.jpg,invalid_status`,
			expectedCount: 0,
			wantErr:       true,
			errContains:   "invalid status value",
		},
		{
			name: "missing required field (title)",
			csvContent: `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,,Description,100.00,image.jpg,available`,
			expectedCount: 0,
			wantErr:       true,
			errContains:   "title is required",
		},
		{
			name: "title too long",
			csvContent: `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,` + string(make([]byte, 201)) + `,Description,100.00,image.jpg,available`,
			expectedCount: 0,
			wantErr:       true,
			errContains:   "title is required",
		},
		{
			name: "duplicate ID",
			csvContent: `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,Item 1,Description 1,100.00,image1.jpg,available
550e8400-e29b-41d4-a716-446655440001,Item 2,Description 2,200.00,image2.jpg,available`,
			expectedCount: 1,
			wantErr:       true,
			errContains:   "already exists",
		},
		{
			name: "wrong number of fields",
			csvContent: `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,Test Item,Description,100.00`,
			expectedCount: 0,
			wantErr:       true,
			errContains:   "wrong number of fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary CSV file
			tmpDir := t.TempDir()
			csvPath := filepath.Join(tmpDir, "test.csv")

			if err := os.WriteFile(csvPath, []byte(tt.csvContent), 0644); err != nil {
				t.Fatalf("Failed to create test CSV file: %v", err)
			}

			// Create repository and loader
			repo := repository.NewMemoryRepository()
			loader := NewCSVLoader(repo)

			// Load CSV
			ctx := context.Background()
			count, err := loader.LoadFromFile(ctx, csvPath)

			// Check error
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error to contain '%s', got '%v'", tt.errContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			// Check count
			if count != tt.expectedCount {
				t.Errorf("expected count %d, got %d", tt.expectedCount, count)
			}

			// Verify items in repository
			items, err := repo.List(ctx)
			if err != nil {
				t.Fatalf("Failed to list items: %v", err)
			}
			if len(items) != tt.expectedCount {
				t.Errorf("expected %d items in repository, got %d", tt.expectedCount, len(items))
			}
		})
	}
}

func TestCSVLoader_LoadFromFile_FileNotFound(t *testing.T) {
	repo := repository.NewMemoryRepository()
	loader := NewCSVLoader(repo)

	ctx := context.Background()
	_, err := loader.LoadFromFile(ctx, "nonexistent.csv")

	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
	if !contains(err.Error(), "failed to open CSV file") {
		t.Errorf("expected 'failed to open CSV file' error, got: %v", err)
	}
}

func TestCSVLoader_LoadFromFile_ValidatesAllFields(t *testing.T) {
	csvContent := `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,Vintage Jacket,Brown leather jacket,150.00,jacket.jpg,available`

	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	repo := repository.NewMemoryRepository()
	loader := NewCSVLoader(repo)

	ctx := context.Background()
	count, err := loader.LoadFromFile(ctx, csvPath)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 item, got %d", count)
	}

	// Verify item details
	item, err := repo.GetByID(ctx, "550e8400-e29b-41d4-a716-446655440001")
	if err != nil {
		t.Fatalf("Failed to get item: %v", err)
	}

	if item.Title != "Vintage Jacket" {
		t.Errorf("expected title 'Vintage Jacket', got '%s'", item.Title)
	}
	if item.Description != "Brown leather jacket" {
		t.Errorf("expected description 'Brown leather jacket', got '%s'", item.Description)
	}
	if item.Price != 150.00 {
		t.Errorf("expected price 150.00, got %.2f", item.Price)
	}
	if item.Image != "jacket.jpg" {
		t.Errorf("expected image 'jacket.jpg', got '%s'", item.Image)
	}
	if item.Status != models.StatusAvailable {
		t.Errorf("expected status 'available', got '%s'", item.Status)
	}
	if item.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if item.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestValidateHeader(t *testing.T) {
	tests := []struct {
		name     string
		actual   []string
		expected []string
		want     bool
	}{
		{
			name:     "valid header",
			actual:   []string{"id", "title", "description", "price", "image", "status"},
			expected: []string{"id", "title", "description", "price", "image", "status"},
			want:     true,
		},
		{
			name:     "wrong order",
			actual:   []string{"title", "id", "description", "price", "image", "status"},
			expected: []string{"id", "title", "description", "price", "image", "status"},
			want:     false,
		},
		{
			name:     "missing field",
			actual:   []string{"id", "title", "description", "price", "image"},
			expected: []string{"id", "title", "description", "price", "image", "status"},
			want:     false,
		},
		{
			name:     "extra field",
			actual:   []string{"id", "title", "description", "price", "image", "status", "extra"},
			expected: []string{"id", "title", "description", "price", "image", "status"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateHeader(tt.actual, tt.expected)
			if got != tt.want {
				t.Errorf("validateHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRecord(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		record  []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid record",
			record:  []string{"id1", "Title", "Description", "100.00", "image.jpg", "available"},
			wantErr: false,
		},
		{
			name:    "wrong number of fields",
			record:  []string{"id1", "Title", "Description"},
			wantErr: true,
			errMsg:  "invalid record length",
		},
		{
			name:    "invalid price",
			record:  []string{"id1", "Title", "Description", "not_a_number", "image.jpg", "available"},
			wantErr: true,
			errMsg:  "invalid price value",
		},
		{
			name:    "invalid status",
			record:  []string{"id1", "Title", "Description", "100.00", "image.jpg", "invalid"},
			wantErr: true,
			errMsg:  "invalid status value",
		},
		{
			name:    "empty title",
			record:  []string{"id1", "", "Description", "100.00", "image.jpg", "available"},
			wantErr: true,
			errMsg:  "item validation failed",
		},
		{
			name:    "negative price",
			record:  []string{"id1", "Title", "Description", "-10.00", "image.jpg", "available"},
			wantErr: true,
			errMsg:  "item validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := parseRecord(tt.record, now)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error to contain '%s', got '%v'", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if item == nil {
					t.Fatal("expected item to be non-nil")
				}
				if item.ID != tt.record[0] {
					t.Errorf("expected ID '%s', got '%s'", tt.record[0], item.ID)
				}
				if item.Title != tt.record[1] {
					t.Errorf("expected Title '%s', got '%s'", tt.record[1], item.Title)
				}
			}
		})
	}
}

func TestNewCSVLoader(t *testing.T) {
	repo := repository.NewMemoryRepository()
	loader := NewCSVLoader(repo)

	if loader == nil {
		t.Fatal("expected loader to be non-nil")
	}
	if loader.repo == nil {
		t.Error("expected loader.repo to be non-nil")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Test helper to verify repository errors are properly wrapped
func TestCSVLoader_LoadFromFile_RepositoryError(t *testing.T) {
	csvContent := `id,title,description,price,image,status
550e8400-e29b-41d4-a716-446655440001,Test Item,Description,100.00,image.jpg,available`

	tmpDir := t.TempDir()
	csvPath := filepath.Join(tmpDir, "test.csv")
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	// Use a mock repository that returns an error
	mockRepo := &errorRepository{err: errors.New("repository error")}
	loader := NewCSVLoader(mockRepo)

	ctx := context.Background()
	count, err := loader.LoadFromFile(ctx, csvPath)

	if err == nil {
		t.Error("expected error from repository, got nil")
	}
	if count != 0 {
		t.Errorf("expected count 0 when repository fails, got %d", count)
	}
	if !contains(err.Error(), "failed to create item") {
		t.Errorf("expected 'failed to create item' error, got: %v", err)
	}
}

// Mock repository that always returns an error
type errorRepository struct {
	err error
}

func (r *errorRepository) Create(ctx context.Context, item *models.Item) error {
	return r.err
}

func (r *errorRepository) GetByID(ctx context.Context, id string) (*models.Item, error) {
	return nil, r.err
}

func (r *errorRepository) List(ctx context.Context) ([]*models.Item, error) {
	return nil, r.err
}

func (r *errorRepository) Update(ctx context.Context, item *models.Item) error {
	return r.err
}

func (r *errorRepository) Delete(ctx context.Context, id string) error {
	return r.err
}
