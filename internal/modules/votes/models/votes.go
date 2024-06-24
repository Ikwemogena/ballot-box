package models

import "time"

type Vote struct {
	ID           string       `json:"id"`
    ElectionID   string       `json:"election_id"`
	VoterID      string       `json:"voter_id"`
	ContestantID string       `json:"contestant_id"`
	CreatedAt    time.Time `json:"created_at"`
}
