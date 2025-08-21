package models

import (
	"encoding/json"
	"time"
)

type TaskState int

const (
	Pending TaskState = iota
	Downloading
	Updating
	Done
	Error
)

// TODO Delete
type TaskLog struct {
	Id        int              `json:"id" db:"id"`
	StartTime time.Time        `json:"startTime" db:"starttime"`
	EndTime   *time.Time       `json:"endTime,omitempty" db:"endtime"`     // Nullable
	Status    int              `json:"status" db:"status"`                 // Expected to be 0–4
	OutputLog *json.RawMessage `json:"outputLog,omitempty" db:"outputlog"` // JSONB field
}

type Song struct {
	Id            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	SourceId      string    `json:"source_id" db:"sourceid"`
	ThumbnailPath string    `json:"thumbnail_path" db:"thumbnailpath"`
	Path          string    `json:"path,omitempty" db:"path"`
	Duration      int       `json:"duration,omitempty" db:"duration"`
	CreatedAt     time.Time `json:"created_at" db:"createdat"`
	UpdatedAt     time.Time `json:"updated_at" db:"updatedat"`
}

type Playlist struct {
	Id            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   string    `json:"description,omitempty" db:"description"`
	ThumbnailPath string    `json:"thumbnailPath" db:"thumbnailpath"`
	CreationDate  time.Time `json:"creationDate" db:"creationdate"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updatedat"`
	IsPublic      bool      `json:"isPublic" db:"ispublic"`
}

type PlaylistSong struct {
	SongId     int       `json:"songId" db:"songid"`
	PlaylistId int       `json:"playlistId" db:"playlistid"`
	Position   int       `json:"position,omitempty" db:"position"`
	AddedAt    time.Time `json:"addedAt" db:"addedat"`
}

type MigrationFile struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"filename" db:"filename"`
	AppliedOn time.Time `json:"appliedon" db:"appliedon"`
}

type ParentTaskLog struct {
	Id      int       `db:"id" json:"id"`
	Url     string    `db:"url" json:"url"`
	AddTime time.Time `db:"add_time" json:"add_time"`
}

type ChildTaskLog struct {
	Id        int             `db:"id" json:"id"`
	ParentId  int             `db:"parent_id" json:"parent_id"`
	StartTime time.Time       `db:"start_time" json:"start_time"`
	EndTime   *time.Time      `db:"end_time" json:"end_time,omitempty"`
	Status    int             `db:"status" json:"status"`
	OutputLog json.RawMessage `db:"output_log" json:"output_log"`
}
