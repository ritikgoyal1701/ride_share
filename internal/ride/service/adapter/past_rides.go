package adapter

import (
	"rideShare/internal/domain/models"
	responses2 "rideShare/internal/ride/service/responses"
)

func GetPastRides(rides []models.Ride) (pastRides []responses2.GetRide) {
	pastRides = make([]responses2.GetRide, 0)
	for _, ride := range rides {
		pastRides = append(pastRides, responses2.GetRide{
			ID: ride.ID.Hex(),
			Rider: responses2.User{
				ID:    ride.Rider.ID,
				Name:  ride.Rider.Name,
				Email: ride.Rider.Email,
			},
			Driver: responses2.User{
				ID:    ride.Driver.ID,
				Name:  ride.Driver.Name,
				Email: ride.Driver.Email,
			},
			StartLocation: getLocationResp(ride.StartLocation.Coordinates),
			DropLocation:  getLocationResp(ride.DropLocation.Coordinates),
			Status:        models.RideStatusToString[ride.Status],
			Price:         ride.Price,
			Distance:      ride.Distance,
		})
	}
	return
}
