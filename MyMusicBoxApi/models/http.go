package models

import "mime/multipart"

type DownloadRequestModel struct {
	Url string `json:"url"`
}

type CreatePlaylistModel struct {
	Name        string                `form:"playlistName"`
	Image       *multipart.FileHeader `form:"backgroundImage"`
	IsPublic    string                `form:"publicPlaylist"`
	Description string                `form:"playlistDescription"`
}

type ApiResponseModel struct {
	Data    any
	Message string
}

func ErrorResponse(data any) ApiResponseModel {
	return ApiResponseModel{
		Data:    data,
		Message: "An error occurred",
	}
}
func OkResponse(data any, message string) ApiResponseModel {
	return ApiResponseModel{
		Data:    data,
		Message: message,
	}
}
