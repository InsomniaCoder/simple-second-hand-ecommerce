package models

import (
	"errors"
	"time"
)

// ItemStatus represents the current status of an item
type ItemStatus string

const (
	StatusAvailable ItemStatus = "available"
	StatusBooked    ItemStatus = "booked"
	StatusSold      ItemStatus = "sold"
)

// Item represents an item for sale in the e-commerce system
type Item struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Image       string     `json:"image"`
	Status      ItemStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Validation errors
var (
	ErrInvalidTitle  = errors.New("title is required and must not exceed 200 characters")
	ErrInvalidPrice  = errors.New("price must be greater than 0")
	ErrInvalidStatus = errors.New("invalid status value")
)

// IsValidStatus checks if the given status is valid
func IsValidStatus(status ItemStatus) bool {
	return status == StatusAvailable || status == StatusBooked || status == StatusSold
}

// Validate performs basic validation on the item
func (i *Item) Validate() error {
	if i.Title == "" || len(i.Title) > 200 {
		return ErrInvalidTitle
	}
	if i.Price <= 0 {
		return ErrInvalidPrice
	}
	if !IsValidStatus(i.Status) {
		return ErrInvalidStatus
	}
	return nil
}
