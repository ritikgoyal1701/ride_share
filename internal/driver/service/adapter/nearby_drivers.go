package adapter

import (
	"rideShare/internal/domain/models"
	"rideShare/internal/driver/service/responses"
)

func GetNearbyDrivers(drivers []models.Driver) (driversResp []responses.Driver) {
	driversResp = make([]responses.Driver, 0)
	for _, driver := range drivers {
		driverResp := responses.Driver{
			User: responses.User{
				ID:    driver.ID.Hex(),
				Name:  driver.Name,
				Email: driver.Email,
			},
		}

		if len(driver.Location.Coordinates) == 2 {
			driverResp.Location = responses.Location{
				XCoordinate: driver.Location.Coordinates[0],
				YCoordinate: driver.Location.Coordinates[1],
			}
		}

		driversResp = append(driversResp, driverResp)
	}

	return
}
