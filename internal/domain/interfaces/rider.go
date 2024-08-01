package interfaces

import (
	"context"
	"github.com/gin-gonic/gin"
	"rideShare/internal/controllers/driver/requests"
	requests2 "rideShare/internal/controllers/rider/requests"
	"rideShare/internal/domain/models"
	"rideShare/internal/rider/service/responses"
	mongo2 "rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
)

type (
	RiderController interface {
		CreateRider(ctx *gin.Context)
	}

	RiderService interface {
		CreateRider(ctx context.Context, req requests2.CreateRiderRequest) (cusErr error2.CustomError)
		Login(ctx context.Context, req requests.LoginRequest) (resp responses.LoginResp, cusErr error2.CustomError)
	}

	RiderRepository interface {
		GetRidersCount(ctx context.Context, filters map[string]mongo2.QueryFilter) (count int64, cusErr error2.CustomError)
		GetRiders(ctx context.Context, filters map[string]mongo2.QueryFilter, fields map[string]interface{}) (riders []models.Rider, cusErr error2.CustomError)
		GetRider(ctx context.Context, filters map[string]mongo2.QueryFilter) (rider *models.Rider, cusErr error2.CustomError)
		CreateRider(ctx context.Context, rider *models.Rider) (cusErr error2.CustomError)
		UpdateRider(ctx context.Context, filters map[string]mongo2.QueryFilter, updates map[string]interface{}) (cusErr error2.CustomError)
	}
)
