package interfaces

import (
	"context"
	"github.com/gin-gonic/gin"
	"rideShare/internal/controllers/driver/requests"
	"rideShare/internal/domain/models"
	mongo2 "rideShare/pkg/db/mongo"
	error2 "rideShare/pkg/error"
)

type (
	DriverController interface {
		CreateDriver(ctx *gin.Context)
	}

	DriverService interface {
		CreateDriver(ctx context.Context, req requests.CreateDriverRequest) (cusErr error2.CustomError)
	}

	DriverRepository interface {
		CreateDriver(ctx context.Context, driver *models.Driver) (cusErr error2.CustomError)
		GetDriver(ctx context.Context, filters map[string]mongo2.QueryFilter) (driver *models.Driver, cusErr error2.CustomError)
		GetDrivers(ctx context.Context, filters map[string]mongo2.QueryFilter, fields map[string]interface{}) (drivers []models.Driver, cusErr error2.CustomError)
		GetDriversCount(ctx context.Context, filters map[string]mongo2.QueryFilter) (count int64, cusErr error2.CustomError)
	}
)