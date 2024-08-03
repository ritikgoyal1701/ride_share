package ride

import (
	"github.com/google/wire"
	"rideShare/internal/domain/interfaces"
	"rideShare/internal/driver"
	driverService "rideShare/internal/driver/service"
	"rideShare/internal/ride"
	rideService "rideShare/internal/ride/service"
	"rideShare/internal/rider"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	NewController,
	driverService.NewService,
	driver.NewRepository,
	rideService.NewService,
	ride.NewRepository,
	rider.NewRepository,

	wire.Bind(new(interfaces.RideController), new(*Controller)),
	wire.Bind(new(interfaces.DriverService), new(*driverService.Service)),
	wire.Bind(new(interfaces.DriverRepository), new(*driver.Repository)),
	wire.Bind(new(interfaces.RideService), new(*rideService.Service)),
	wire.Bind(new(interfaces.RideRepository), new(*ride.Repository)),
	wire.Bind(new(interfaces.RiderRepository), new(*rider.Repository)),
)
