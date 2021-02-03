package models

import "time"

type Check struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=30"`
	URL      string `json:"url" validate:"required,url"`
	Interval uint64 `json:"interval" validate:"required,min=1"`
	// Statuses []Status  `json:"-"`
	Email    string `json:"email" validate:"required,email"`
	Active   bool   `json:"active" validate:"required"`
}

type CheckUpdate struct {
	URL      *string `json:"url" validate:"url"`
	Name     *string `json:"name" validate:"min=3,max=30"`
	Interval *uint64 `json:"interval" validate:"min=1"`
	Email    *string `json:"email" validate:"email"`
	Active   *bool   `json:"active"`
}

type Status struct {
	ID      string `json:"-"`
	CheckID string `json:"-" db:"check_id"`
	Content string `json:"content"`// TODO byte array maybe
	Date    time.Time `json:"date"`
}