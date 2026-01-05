package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
)

// NewTestItem creates a test item with default values
func NewTestItem(overrides ...func(*models.Item)) *models.Item {
	item := &models.Item{
		ID:          "test-123",
		Title:       "Test Item",
		Description: "Test Description",
		Price:       10.99,
		Image:       "https://example.com/image.jpg",
		Status:      models.StatusAvailable,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, override := range overrides {
		override(item)
	}

	return item
}

// NewTestItemMinimal creates a minimal test item
func NewTestItemMinimal() *models.Item {
	return &models.Item{
		ID:     "test-minimal",
		Title:  "Minimal",
		Price:  1.00,
		Status: models.StatusAvailable,
	}
}

// AssertEqual checks if two values are equal
func AssertEqual(t *testing.T, got, want interface{}, message string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", message, got, want)
	}
}

// AssertNotNil checks if a value is not nil
func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Errorf("%s: expected non-nil value", message)
	}
}

// AssertNil checks if a value is nil
func AssertNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value != nil {
		t.Errorf("%s: expected nil, got %v", message, value)
	}
}

// AssertNoError checks if error is nil
func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: unexpected error: %v", message, err)
	}
}

// AssertError checks if error is not nil
func AssertError(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error, got nil", message)
	}
}

// AssertHTTPStatus checks HTTP response status code
func AssertHTTPStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("HTTP status: got %d, want %d", got, want)
	}
}

// AssertJSONResponse checks if response contains valid JSON
func AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	t.Helper()
	if err := json.NewDecoder(w.Body).Decode(target); err != nil {
		t.Errorf("failed to decode JSON response: %v", err)
	}
}

// MakeJSONRequest creates an HTTP request with JSON body
func MakeJSONRequest(method, url string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// ItemsEqual checks if two items have the same field values (ignoring timestamps)
func ItemsEqual(t *testing.T, got, want *models.Item) {
	t.Helper()
	if got.ID != want.ID {
		t.Errorf("ID: got %v, want %v", got.ID, want.ID)
	}
	if got.Title != want.Title {
		t.Errorf("Title: got %v, want %v", got.Title, want.Title)
	}
	if got.Description != want.Description {
		t.Errorf("Description: got %v, want %v", got.Description, want.Description)
	}
	if got.Price != want.Price {
		t.Errorf("Price: got %v, want %v", got.Price, want.Price)
	}
	if got.Image != want.Image {
		t.Errorf("Image: got %v, want %v", got.Image, want.Image)
	}
	if got.Status != want.Status {
		t.Errorf("Status: got %v, want %v", got.Status, want.Status)
	}
}

// ValidItemFixtures returns a slice of valid test items
func ValidItemFixtures() []*models.Item {
	return []*models.Item{
		{
			ID:          "item-1",
			Title:       "MacBook Pro",
			Description: "2021 model, M1 Max",
			Price:       1999.99,
			Image:       "https://example.com/macbook.jpg",
			Status:      models.StatusAvailable,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "item-2",
			Title:       "Desk Chair",
			Description: "Ergonomic office chair",
			Price:       150.00,
			Image:       "https://example.com/chair.jpg",
			Status:      models.StatusBooked,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "item-3",
			Title:       "Monitor",
			Description: "27-inch 4K display",
			Price:       399.99,
			Image:       "https://example.com/monitor.jpg",
			Status:      models.StatusSold,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

// InvalidItemFixtures returns items that should fail validation
func InvalidItemFixtures() []*models.Item {
	return []*models.Item{
		{
			// Empty title
			ID:    "invalid-1",
			Title: "",
			Price: 10.00,
		},
		{
			// Negative price
			ID:    "invalid-2",
			Title: "Item",
			Price: -10.00,
		},
		{
			// Invalid status
			ID:     "invalid-3",
			Title:  "Item",
			Price:  10.00,
			Status: "invalid-status",
		},
		{
			// Title too long
			ID:    "invalid-4",
			Title: string(make([]byte, 201)),
			Price: 10.00,
		},
	}
}
