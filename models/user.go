package models

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"user_password"`
	Role         string `json:"user_role"`
	Access_token string `json:"access_token"`
}
