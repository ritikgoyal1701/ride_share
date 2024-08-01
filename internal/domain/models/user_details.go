package models

type UserDetails struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Title Title  `json:"title"`
}
