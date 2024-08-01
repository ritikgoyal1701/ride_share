package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"rideShare/constants"
	"rideShare/internal/domain/models"
	jwt2 "rideShare/pkg/jwt"
	"strings"
)

func AuthenticateJWT(title models.Title) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken, err := extractToken(c)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, isValid, err := jwt2.ValidateToken(bearerToken)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !isValid {
			log.Printf("Token is invalid :: %v\n", bearerToken)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if len(title) > 0 && claims.Title != title {
			log.Printf("Invalid hit")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(constants.ID, claims.ID)
		c.Set(constants.Email, claims.Email)
		c.Set(constants.Title, claims.Title)
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	bearerToken := c.Request.Header.Get(constants.Authorization)
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1], nil
	}
	return "", errors.New("no auth token found")
}
