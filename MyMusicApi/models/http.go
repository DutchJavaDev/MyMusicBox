package models

type UrlRequest struct {
	Url string `json:"url"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
