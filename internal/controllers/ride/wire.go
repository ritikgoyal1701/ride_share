//go:build wireinject
// +build wireinject

package ride

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func RideControllerWire(client *mongo.Database) (*Controller, error) {
	panic(wire.Build(ProviderSet))
}
