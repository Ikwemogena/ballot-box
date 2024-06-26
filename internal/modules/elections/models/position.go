package models

import "time"

type Position struct {
	ID          string       `json:"id"`
	Name		string    `json:"name"`
	ElectionID       string    `json:"election_id"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
	Contestants []Contestant `json:"contestants"`
}