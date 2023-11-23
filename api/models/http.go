package models

type Request struct {
	URL    string `json:"url"`
	Custom string `json:"custom,omitempty"`
}

type Response struct {
	URL     string `json:"url"`
	NewURL  string `json:"new_url,omitempty"`
	Message string `json:"message,omitempty"`
}

type ResponseData struct {
	Message string            `json:"message,omitempty"`
	Error   int               `json:"error,omitempty"`
	Data    map[string]string `json:"data,omitempty"`
}
