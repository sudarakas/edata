package types

import (
	"context"
	"time"
)

// Store
type SubscriptionStore interface {
	CreateSubscription(ctx context.Context, subscription Subscription) (Subscription, error)
	// GetAllSubscriptions(ctx context.Context) ([]Subscription, error)
	// GetSubscriptionByID(ctx context.Context, id string) (Subscription, error)
	GetSubscriptionByCode(ctx context.Context, code string) (Subscription, error)
	UpdateSubscription(ctx context.Context, subscription Subscription) (Subscription, error)
	DeleteSubscription(ctx context.Context, userID string) error
	// FilterSubscriptions(ctx context.Context, filter SubscriptionFilter) ([]Subscription, error)
}

// Models
type Subscription struct {
	SubscriptionID string    `json:"id"`
	Code           string    `json:"code"`
	Plan           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	DataLimit      int       `json:"data_limit"`
	Validity       int       `json:"validity_period"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}

// Filter
type SubscriptionFilter struct {
	Plan        string `json:"plan"`
	MinPrice    int    `json:"min_price"`
	MaxPrice    int    `json:"max_price"`
	MinData     int    `json:"min_data"`
	MaxData     int    `json:"max_data"`
	MinValidity int    `json:"min_validity"`
	MaxValidity int    `json:"max_validity"`
}

// Payloads
type CreateSubscriptionPayLoad struct {
	Code        string  `json:"code" validate:"required,min=2,max=100"`
	Plan        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description" validate:"required,min=2,max=1000"`
	Price       float64 `json:"price" validate:"required,min=0"`
	DataLimit   int     `json:"data_limit" validate:"required,min=0"`
	Validity    int     `json:"validity_period" validate:"required,min=0"`
}

type UpdateSubscriptionPayLoad struct {
	Code        string  `json:"code" validate:"required,min=2,max=100,omitempty"`
	Plan        string  `json:"name" validate:"required,min=2,max=100,omitempty"`
	Description string  `json:"description" validate:"required,min=2,max=1000,omitempty"`
	Price       float64 `json:"price" validate:"required,min=0,omitempty"`
	DataLimit   int     `json:"data_limit" validate:"required,min=0,omitempty"`
	Validity    int     `json:"validity_period" validate:"required,min=0,omitempty"`
}

type DeleteSubscriptionPayLoad struct {
	SubscriptionID string `json:"id" validate:"required"`
}

type FilterSubscriptionPayLoad struct {
	Plan        string `json:"plan"`
	MinPrice    int    `json:"min_price"`
	MaxPrice    int    `json:"max_price"`
	MinData     int    `json:"min_data"`
	MaxData     int    `json:"max_data"`
	MinValidity int    `json:"min_validity"`
	MaxValidity int    `json:"max_validity"`
}
