package utils

import (
	"golang.org/x/net/context"
	"net/http"
	"rideShare/constants"
	"rideShare/internal/domain/models"
	error2 "rideShare/pkg/error"
)

func GetUserDetails(ctx context.Context) (userDetails models.UserDetails, cusErr error2.CustomError) {
	id, ok := ctx.Value(constants.ID).(string)
	if !ok {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, "user id not found in ctx")
		return
	}

	email, ok := ctx.Value(constants.Email).(string)
	if !ok {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, "user email not found in ctx")
		return
	}

	title, ok := ctx.Value(constants.Title).(string)
	if !ok {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, "user title not found in ctx")
		return
	}

	modelTitle, ok := models.StringToTitle[title]
	if !ok {
		cusErr = error2.NewCustomError(http.StatusInternalServerError, "Invalid title")
		return
	}

	userDetails = models.UserDetails{
		ID:    id,
		Email: email,
		Title: modelTitle,
	}
	return
}
