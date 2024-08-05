package models

import "github.com/golang-jwt/jwt"

type (
	Title string
)

const (
	TitleDriver Title = "driver"
	TitleRider  Title = "rider"
)

var TitleToString = map[Title]string{
	TitleDriver: "driver",
	TitleRider:  "rider",
}

var StringToTitle = map[string]Title{
	"driver": TitleDriver,
	"rider":  TitleRider,
}

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Title Title  `json:"title"`
	jwt.StandardClaims
}

func (t Title) IsDriver() bool {
	return t == TitleDriver
}

func (t Title) IsRider() bool {
	return t == TitleRider
}
