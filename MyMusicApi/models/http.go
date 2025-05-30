package models

type UrlRequest struct {
	Url string `json:"url"`
}

type DownloadRequest struct {
	Url        string `json:"url"`
	Playlist   string `json:"playlist"`   // name of playlist to add to
	IsPlaylist bool   `json:"isPlaylist"` // yt playlist or single video
}

type Response struct {
	Data    any
	Message string
}

func ErrorResponse(data any) Response {
	return Response{
		Data:    data,
		Message: "An error occurred",
	}
}
func OkResponse(data any, message string) Response {
	return Response{
		Data:    data,
		Message: message,
	}
}
