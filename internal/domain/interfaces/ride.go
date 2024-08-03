package interfaces

import (
	"context"
	"github.com/gin-gonic/gin"
	requests2 "rideShare/internal/controllers/ride/requests"
	"rideShare/internal/domain/models"
	"rideShare/internal/driver/service/responses"
	responses2 "rideShare/internal/ride/service/responses"
	mongo2 "rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
)

type (
	RideController interface {
		GetRidePrice(ctx *gin.Context)
		CreateRide(ctx *gin.Context)
		GetRides(ctx *gin.Context)
	}

	RideService interface {
		CreateRide(ctx context.Context, userDetails models.UserDetails, req requests2.CreateRideRequest) (cusErr error2.CustomError)
		GetRidePrice(ctx context.Context, userDetails models.UserDetails, req requests2.PriceRequest) (resp responses.PriceResponse, cusErr error2.CustomError)
		GetRides(ctx context.Context, userDetails models.UserDetails, req requests2.Location) (resp []responses2.GetRides, cusErr error2.CustomError)
	}

	RideRepository interface {
		CreateRide(ctx context.Context, ride *models.Ride) (cusErr error2.CustomError)
		GetRide(ctx context.Context, filters map[string]mongo2.QueryFilter) (ride *models.Ride, cusErr error2.CustomError)
		GetRides(ctx context.Context, filters map[string]mongo2.QueryFilter, fields map[string]interface{}) (riders []models.Ride, cusErr error2.CustomError)
		GetRidesCount(ctx context.Context, filters map[string]mongo2.QueryFilter) (count int64, cusErr error2.CustomError)
		UpdateRide(ctx context.Context, filters map[string]mongo2.QueryFilter, updates map[string]interface{}) (cusErr error2.CustomError)
	}
)
