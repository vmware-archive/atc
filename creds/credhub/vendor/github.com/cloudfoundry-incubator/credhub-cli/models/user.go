package models

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
