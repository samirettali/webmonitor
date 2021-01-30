package models

type Job struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Interval uint64 `json:"interval"`
	State    string `json:"-"`
	Email    string `json:"email"`
}

type JobUpdate struct {
	URL      *string `json:"url"`
	Interval *uint64 `json:"interval"`
	Email    *string `json:"email"`
	State    *string `json:"-"`
}
