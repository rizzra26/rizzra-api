package models

import "time"

type StackItem struct {
	ID          string     `json:"id"`
	CategoryID  string     `json:"category_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	OrderIndex  int        `json:"order_index"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
