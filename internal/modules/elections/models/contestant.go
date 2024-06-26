package models

import "time"

type Contestant struct {
	ID         string       `json:"id"`
	ElectionID string       `json:"election_id"`
	Name       string    `json:"name"`
	PositionID string   `json:"position_id"`
	CreatedAt  time.Time `json:"created_at"`
}