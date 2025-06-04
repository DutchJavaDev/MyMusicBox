package models

import (
	"encoding/json"
	"time"
)

const (
	LogTypeInfo = iota
	LogTypeWarning
	LogTypeError
)

type Log struct {
	Id        int             `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Message   string          `json:"message"`
	Type      int             `json:"type"`    // 0 = info, 1 = warning, 2 = error
	Context   json.RawMessage `json:"context"` // use a struct instead if structure is known
}
