package loader

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/models"
	"github.com/InsomniaCoder/simple-second-hand-ecommerce/internal/repository"
)

// CSVLoader handles loading items from CSV files
type CSVLoader struct {
	repo repository.ItemRepository
}

// NewCSVLoader creates a new CSV loader
func NewCSVLoader(repo repository.ItemRepository) *CSVLoader {
	return &CSVLoader{
		repo: repo,
	}
}

// LoadFromFile loads items from a CSV file and stores them in the repository
func (l *CSVLoader) LoadFromFile(ctx context.Context, filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return 0, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate header
	expectedHeader := []string{"id", "title", "description", "price", "image", "status"}
	if !validateHeader(header, expectedHeader) {
		return 0, fmt.Errorf("invalid CSV header: expected %v, got %v", expectedHeader, header)
	}

	count := 0
	now := time.Now()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, fmt.Errorf("failed to read CSV record at line %d: %w", count+2, err)
		}

		item, err := parseRecord(record, now)
		if err != nil {
			return count, fmt.Errorf("failed to parse CSV record at line %d: %w", count+2, err)
		}

		if err := l.repo.Create(ctx, item); err != nil {
			return count, fmt.Errorf("failed to create item %s at line %d: %w", item.ID, count+2, err)
		}

		count++
	}

	return count, nil
}

// validateHeader checks if the CSV header matches the expected format
func validateHeader(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if actual[i] != expected[i] {
			return false
		}
	}
	return true
}

// parseRecord converts a CSV record into an Item
func parseRecord(record []string, timestamp time.Time) (*models.Item, error) {
	if len(record) != 6 {
		return nil, fmt.Errorf("invalid record length: expected 6 fields, got %d", len(record))
	}

	price, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price value '%s': %w", record[3], err)
	}

	status := models.ItemStatus(record[5])
	if !models.IsValidStatus(status) {
		return nil, fmt.Errorf("invalid status value: %s", record[5])
	}

	item := &models.Item{
		ID:          record[0],
		Title:       record[1],
		Description: record[2],
		Price:       price,
		Image:       record[4],
		Status:      status,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	if err := item.Validate(); err != nil {
		return nil, fmt.Errorf("item validation failed: %w", err)
	}

	return item, nil
}
