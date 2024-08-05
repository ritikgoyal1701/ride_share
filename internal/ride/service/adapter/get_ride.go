package adapter

import (
	"rideShare/internal/domain/models"
	responses2 "rideShare/internal/ride/service/responses"
)

func GetRideDetails(ride *models.Ride, drivers []models.Driver) (rideResp responses2.GetRide) {
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

	if len(drivers) > 0 {
		rideResp.Driver = responses2.User{
			ID:    drivers[0].ID.Hex(),
			Name:  drivers[0].Name,
			Email: drivers[0].Email,
		}

		rideResp.DriverLocation = getLocationResp(drivers[0].Location.Coordinates)
	}

	return
}

func getLocationResp(coordinates []float64) responses2.Location {
	return responses2.Location{
		XCoordinate: coordinates[0],
		YCoordinate: coordinates[1],
	}
}
