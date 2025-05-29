package db

import "time"

type Song struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	SourceURL string    `json:"source_url" db:"sourceurl"`
	Path      *string   `json:"path,omitempty" db:"path"`
	Duration  *int      `json:"duration,omitempty" db:"duration"`
	CreatedAt time.Time `json:"created_at" db:"createdat"`
	UpdatedAt time.Time `json:"updated_at" db:"updatedat"`
}
