package adapter

import (
	"rideShare/internal/domain/models"
	responses2 "rideShare/internal/ride/service/responses"
)

func GetRideDetails(ride *models.Ride, driver *models.Driver) (rideResp responses2.GetRide) {
	rideResp = responses2.GetRide{
		ID: ride.ID.Hex(),
		Rider: responses2.User{
			ID:    ride.Rider.ID,
			Name:  ride.Rider.Name,
			Email: ride.Rider.Email,
		},
		StartLocation: getLocationResp(ride.StartLocation.Coordinates),
		DropLocation:  getLocationResp(ride.DropLocation.Coordinates),
		Status:        models.RideStatusToString[ride.Status],
		Price:         ride.Price,
		Distance:      ride.Distance,
	}

	if driver != nil {
		rideResp.Driver = responses2.User{
			ID:    driver.ID.Hex(),
			Name:  driver.Name,
			Email: driver.Email,
		}

		rideResp.DriverLocation = getLocationResp(driver.Location.Coordinates)
	}

	return
}

func getLocationResp(coordinates []float64) responses2.Location {
	return responses2.Location{
		XCoordinate: coordinates[0],
		YCoordinate: coordinates[1],
	}
}
