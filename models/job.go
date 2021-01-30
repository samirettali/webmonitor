package models

type Job struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Interval uint64 `json:"interval"`
	State    string `json:"state"`
	Email    string `json:"email"`
}
