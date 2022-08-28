package blockchain

import "github.com/golang/protobuf/ptypes/timestamp"

type ValidatorActionResult struct {
	ID          string
	FAK         string
	ValidatorID string
	Results     []ValidatorResult
	Action      Action
	Timestamp   timestamp.Timestamp
}

type ValidatorResult struct {
	ValidatorID string
	Result      bool
}
type Action struct {
	ActionID   string
	ActionType string
	ActionData string
	ResourceID string
}
