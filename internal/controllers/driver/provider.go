package driver

import (
	"github.com/google/wire"
	"rideShare/internal/domain/interfaces"
	"rideShare/internal/driver"
	driverService "rideShare/internal/driver/service"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	NewController,
	driverService.NewService,
	driver.NewRepository,

	wire.Bind(new(interfaces.DriverController), new(*Controller)),
	wire.Bind(new(interfaces.DriverService), new(*driverService.Service)),
	wire.Bind(new(interfaces.DriverRepository), new(*driver.Repository)),
)
