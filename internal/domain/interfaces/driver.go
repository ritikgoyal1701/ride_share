package interfaces

import (
	"context"
	"github.com/gin-gonic/gin"
	"rideShare/internal/controllers/driver/requests"
	"rideShare/internal/domain/models"
	"rideShare/internal/driver/service/responses"
	mongo2 "rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
)

type (
	DriverController interface {
		CreateDriver(ctx *gin.Context)
		Login(ctx *gin.Context)
		Logout(ctx *gin.Context)
		UpdateDriverLocation(ctx *gin.Context)
		GetNearbyDrivers(ctx *gin.Context)
	}

	DriverService interface {
		CreateDriver(ctx context.Context, req requests.CreateDriverRequest) (cusErr error2.CustomError)
		Login(ctx context.Context, req requests.LoginRequest) (resp responses.LoginResp, cusErr error2.CustomError)
		Logout(ctx context.Context, userDetails models.UserDetails) (cusErr error2.CustomError)
		UpdateDriverLocation(ctx context.Context, userDetails models.UserDetails, req requests.LocationUpdate) (cusErr error2.CustomError)
		GetNearbyDrivers(ctx context.Context, userDetails models.UserDetails, req requests.NearByDriversRequest) (driversResp []responses.Driver, cusErr error2.CustomError)
	}

	DriverRepository interface {
		CreateDriver(ctx context.Context, driver *models.Driver) (cusErr error2.CustomError)
		GetDriver(ctx context.Context, filters map[string]mongo2.QueryFilter) (driver *models.Driver, cusErr error2.CustomError)
		GetDrivers(ctx context.Context, filters map[string]mongo2.QueryFilter, fields map[string]interface{}) (drivers []models.Driver, cusErr error2.CustomError)
		GetDriversCount(ctx context.Context, filters map[string]mongo2.QueryFilter) (count int64, cusErr error2.CustomError)
		UpdateDriver(ctx context.Context, filters map[string]mongo2.QueryFilter, updates map[string]interface{}) (cusErr error2.CustomError)
	}
)
