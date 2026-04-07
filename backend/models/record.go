package models

import "time"

type Record struct {
	ID            int       `json:"id"`
	KidID         int       `json:"kid_id"`
	BehaviorID    int       `json:"behavior_id"`
	BehaviorName  string    `json:"behavior_name"`
	Points        int       `json:"points"`
	CreatedAt     time.Time `json:"created_at"`
	Reversed      bool      `json:"reversed"`
	ReversedReason string   `json:"reversed_reason,omitempty"`
	ReversedAt    *time.Time `json:"reversed_at,omitempty"`
}
