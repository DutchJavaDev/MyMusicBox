package db

type Music struct {
	Id   int
	Name string `json:"name"`
	Url  string `json:"url"`
	Path string `json:"path"`
}
