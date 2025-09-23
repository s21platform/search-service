package model

import "time"

type User struct {
	UUID      string    `json:"uuid"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}
