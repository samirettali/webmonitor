package models

type Job struct {
	ID       string `json:"id"`
	Name     string `json:"name" validate:"required,min=3,max=30"`
	URL      string `json:"url" validate:"required,url"`
	Interval uint64 `json:"interval" validate:"required,min=1"`
	State    string `json:"-"`
	Email    string `json:"email" validate:"required,email"`
	Active   bool   `json:"active" validate:"required"`
}

type JobUpdate struct {
	URL      *string `json:"url" validate:"url"`
	Name     *string `json:"name" validate:"min=3,max=30"`
	Interval *uint64 `json:"interval" validate:"min=1"`
	Email    *string `json:"email" validate:"email"`
	State    *string `json:"-"`
	Active   *bool   `json:"active"`
}
