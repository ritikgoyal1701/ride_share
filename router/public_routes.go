package router

import (
	"github.com/gin-gonic/gin"
	"rideShare/internal/controllers/driver"
	"rideShare/internal/controllers/rider"
	"rideShare/pkg/db/mongo"
)

func PublicRoutes(r *gin.Engine) (err error) {
	driverController, err := driver.DriverControllerWire(mongo.Get().Database)
	if err != nil {
		return
	}

	riderController, err := rider.RiderControllerWire(mongo.Get().Database)
	if err != nil {
		return
	}

	drivers := r.Group("/api/v1/drivers")
	{
		drivers.POST("/", driverController.CreateDriver)
	}

	riders := r.Group("/api/v1/riders")
	{
		riders.POST("/", riderController.CreateRider)
	}

	return
}
