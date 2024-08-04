package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"rideShare/pkg/validate"
	"rideShare/router"
)

func Initialize(ctx context.Context, r *gin.Engine) (err error) {
	initializeDB(ctx)
	validate.Set()
	initializeRedis(ctx)
	err = router.PublicRoutes(r)
	if err != nil {
		panic(err)
	}

	return
}
