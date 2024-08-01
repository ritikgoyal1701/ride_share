//go:build wireinject
// +build wireinject

package driver

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func DriverControllerWire(client *mongo.Database) (*Controller, error) {
	panic(wire.Build(ProviderSet))
}
