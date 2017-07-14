package models

type ServerError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
