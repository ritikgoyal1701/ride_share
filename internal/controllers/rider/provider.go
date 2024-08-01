package rider

import (
	"github.com/google/wire"
	"rideShare/internal/domain/interfaces"
	"rideShare/internal/rider"
	riderService "rideShare/internal/rider/service"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	NewController,
	riderService.NewService,
	rider.NewRepository,

	wire.Bind(new(interfaces.RiderController), new(*Controller)),
	wire.Bind(new(interfaces.RiderService), new(*riderService.Service)),
	wire.Bind(new(interfaces.RiderRepository), new(*rider.Repository)),
)
