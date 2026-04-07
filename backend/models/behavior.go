package models

type BehaviorType string

const (
	BehaviorTypeAdd BehaviorType = "add"
	BehaviorTypeSub BehaviorType = "sub"
)

type Behavior struct {
	ID     int          `json:"id"`
	Name   string       `json:"name"`
	Type   BehaviorType `json:"type"`
	Points int          `json:"points"`
}
