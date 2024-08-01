package models

import "github.com/golang-jwt/jwt"

type (
	Title string
)

const (
	TitleDriver Title = "driver"
	TitleRider  Title = "rider"
)

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Title Title  `json:"title"`
	jwt.StandardClaims
}
