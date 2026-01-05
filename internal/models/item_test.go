package models

import (
	"errors"
	"testing"
)

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status ItemStatus
		want   bool
	}{
		{
			name:   "available status",
			status: StatusAvailable,
			want:   true,
		},
		{
			name:   "booked status",
			status: StatusBooked,
			want:   true,
		},
		{
			name:   "sold status",
			status: StatusSold,
			want:   true,
		},
		{
			name:   "empty status",
			status: "",
			want:   false,
		},
		{
			name:   "invalid status",
			status: "invalid",
			want:   false,
		},
		{
			name:   "uppercase status",
			status: "AVAILABLE",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidStatus(tt.status)
			if got != tt.want {
				t.Errorf("IsValidStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    *Item
		wantErr bool
		errType error
	}{
		{
			name: "valid item",
			item: &Item{
				Title:       "Test Item",
				Description: "Test Description",
				Price:       10.99,
				Status:      StatusAvailable,
			},
			wantErr: false,
		},
		{
			name: "valid item with booked status",
			item: &Item{
				Title:  "Booked Item",
				Price:  50.00,
				Status: StatusBooked,
			},
			wantErr: false,
		},
		{
			name: "valid item with sold status",
			item: &Item{
				Title:  "Sold Item",
				Price:  100.00,
				Status: StatusSold,
			},
			wantErr: false,
		},
		{
			name: "valid item with max title length",
			item: &Item{
				Title:  string(make([]byte, 200)),
				Price:  10.99,
				Status: StatusAvailable,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			item: &Item{
				Title:  "",
				Price:  10.99,
				Status: StatusAvailable,
			},
			wantErr: true,
			errType: ErrInvalidTitle,
		},
		{
			name: "title exceeds 200 characters",
			item: &Item{
				Title:  string(make([]byte, 201)),
				Price:  10.99,
				Status: StatusAvailable,
			},
			wantErr: true,
			errType: ErrInvalidTitle,
		},
		{
			name: "zero price",
			item: &Item{
				Title:  "Test Item",
				Price:  0,
				Status: StatusAvailable,
			},
			wantErr: true,
			errType: ErrInvalidPrice,
		},
		{
			name: "negative price",
			item: &Item{
				Title:  "Test Item",
				Price:  -10.99,
				Status: StatusAvailable,
			},
			wantErr: true,
			errType: ErrInvalidPrice,
		},
		{
			name: "invalid status",
			item: &Item{
				Title:  "Test Item",
				Price:  10.99,
				Status: "invalid",
			},
			wantErr: true,
			errType: ErrInvalidStatus,
		},
		{
			name: "empty status",
			item: &Item{
				Title:  "Test Item",
				Price:  10.99,
				Status: "",
			},
			wantErr: true,
			errType: ErrInvalidStatus,
		},
		{
			name: "multiple validation errors - title and price",
			item: &Item{
				Title:  "",
				Price:  -5.00,
				Status: StatusAvailable,
			},
			wantErr: true,
			errType: ErrInvalidTitle, // First error encountered
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()

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

func TestItemStatus_Constants(t *testing.T) {
	// Verify status constants have expected values
	if StatusAvailable != "available" {
		t.Errorf("StatusAvailable = %q, want %q", StatusAvailable, "available")
	}
	if StatusBooked != "booked" {
		t.Errorf("StatusBooked = %q, want %q", StatusBooked, "booked")
	}
	if StatusSold != "sold" {
		t.Errorf("StatusSold = %q, want %q", StatusSold, "sold")
	}
}
