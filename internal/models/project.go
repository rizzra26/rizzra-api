package models

import "time"

type Project struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Tech        []string   `json:"tech"`
	GithubURL   *string    `json:"github_url"`
	DemoURL     *string    `json:"demo_url"`
	CoverURL    *string    `json:"cover_url"`
	OrderIndex  int        `json:"order_index"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type ProjectReorderItem struct {
	ID         string `json:"id" validate:"required,uuid"`
	OrderIndex int    `json:"order_index" validate:"gte=0"`
}
