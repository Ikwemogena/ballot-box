package models

import "time"

type Contestant struct {
	ID         string       `json:"id"`
	ElectionID string       `json:"election_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
}