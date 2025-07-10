package models

type YtdlpJsonResult struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Duration      int    `json:"duration"`
	Playlist      string `json:"playlist"`
	PlaylistIndex int    `json:"index"`
}
