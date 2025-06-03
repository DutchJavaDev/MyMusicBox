package models

type DownloadRequestModel struct {
	Url string `json:"url"`
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
