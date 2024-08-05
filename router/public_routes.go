package router

import (
	"github.com/gin-gonic/gin"
	"rideShare/internal/controllers/driver"
	"rideShare/internal/controllers/ride"
	"rideShare/internal/controllers/rider"
	"rideShare/internal/domain/models"
	"rideShare/pkg/db/mongo"
	"rideShare/pkg/redis"
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

	rideController, err := ride.RideControllerWire(mongo.Get().Database, redis.GetClient().Client)
	if err != nil {
		return
	}

	drivers := r.Group("/api/v1/drivers")
	{
		drivers.POST("/", driverController.CreateDriver)
		drivers.POST("/login", driverController.Login)
		drivers.PUT("/logout", utils.AuthenticateJWT(models.TitleDriver), driverController.Logout)
		drivers.PUT("/location", utils.AuthenticateJWT(models.TitleDriver), driverController.UpdateDriverLocation)
		drivers.GET("/nearby", utils.AuthenticateJWT(models.TitleRider), driverController.GetNearbyDrivers)
	}

	rides := r.Group("/api/v1/rides")
	{
		rides.GET("/price", utils.AuthenticateJWT(models.TitleRider), rideController.GetRidePrice)
		rides.POST("/", utils.AuthenticateJWT(models.TitleRider), rideController.CreateRide)
		rides.PUT("/:ride_id/accept", utils.AuthenticateJWT(models.TitleDriver), rideController.AcceptRide)
		rides.PUT("/:ride_id/cancel", utils.AuthenticateJWT(models.TitleRider), rideController.CancelRide)
		rides.PUT("/:ride_id/verify", utils.AuthenticateJWT(models.TitleDriver), rideController.VerifyRide)
		rides.PUT("/:ride_id/complete", utils.AuthenticateJWT(models.TitleDriver), rideController.CompleteRide)
		rides.GET("/", utils.AuthenticateJWT(models.TitleDriver), rideController.GetRides)
		rides.GET("/:ride_id", utils.AuthenticateJWT(""), rideController.GetRide)
		rides.GET("/past", utils.AuthenticateJWT(""), rideController.GetPastRides)
	}

	riders := r.Group("/api/v1/riders")
	{
		riders.POST("/", riderController.CreateRider)
		riders.POST("/login", riderController.Login)
		riders.PUT("/logout", utils.AuthenticateJWT(models.TitleRider), riderController.Logout)
	}

	return
}
