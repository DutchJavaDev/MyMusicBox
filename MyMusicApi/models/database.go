package models

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

type Playlist struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"` // nullable
	CreationDate time.Time `json:"creationDate"`
	UpdatedAt    time.Time `json:"updatedAt"`
	IsPublic     bool      `json:"isPublic"`
}

type PlaylistSong struct {
	SongID     int       `json:"songId"`
	PlaylistID int       `json:"playlistId"`
	Position   int       `json:"position"`
	AddedAt    time.Time `json:"addedAt"`
}
