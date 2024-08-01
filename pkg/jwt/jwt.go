package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"rideShare/internal/domain/models"
	"time"
)

var jwtKey = []byte("rider_share_desi")

func GenerateToken(email, id string, title models.Title) (string, error) {
	claims := &models.Claims{
		ID:    id,
		Email: email,
		Title: title,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.New().String(),
			ExpiresAt: time.Now().Add(time.Minute * 240).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (claims *models.Claims, isValid bool, err error) {
	claims = &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return
	}

	if !token.Valid {
		return
	}

	isValid = true
	return
}
