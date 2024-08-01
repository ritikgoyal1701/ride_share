package router

import (
	"github.com/gin-gonic/gin"
	"rideShare/internal/controllers/driver"
	"rideShare/internal/controllers/rider"
	"rideShare/internal/domain/models"
	"rideShare/pkg/db/mongo"
	"rideShare/pkg/utils"
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
		drivers.POST("/login", driverController.Login)
		drivers.PUT("/logout", utils.AuthenticateJWT(models.TitleDriver), driverController.Logout)
	}

	riders := r.Group("/api/v1/riders")
	{
		riders.POST("/", riderController.CreateRider)
		riders.POST("/login", riderController.Login)
		riders.PUT("/logout", utils.AuthenticateJWT(models.TitleRider), riderController.Logout)
	}

	return
}
