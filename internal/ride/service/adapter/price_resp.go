package adapter

import (
	"rideShare/internal/controllers/ride/requests"
	"rideShare/internal/domain/models"
	"rideShare/internal/driver/service/responses"
	responses2 "rideShare/internal/ride/service/responses"
	"rideShare/pkg/utils"
)

func GetRidePrice(distance float64, price float64, surge bool, drivers []responses.Driver) responses.PriceResponse {
	return responses.PriceResponse{
		Drivers:  drivers,
		Price:    price,
		Distance: distance,
		IsSurge:  surge,
	}
}

func GetRides(rides []models.Ride, req requests.Location) (ridesResp []responses2.GetRides) {
	ridesResp = make([]responses2.GetRides, 0)
	for _, ride := range rides {
		ridesResp = append(ridesResp, responses2.GetRides{
			ID:            ride.ID.Hex(),
			RiderDistance: utils.CalculateDistance(req.XCoordinate, req.YCoordinate, ride.StartLocation.Coordinates[0], ride.StartLocation.Coordinates[1]),
			RideDistance:  ride.Distance,
			Price:         ride.Price,
		})
	}
	return
}
