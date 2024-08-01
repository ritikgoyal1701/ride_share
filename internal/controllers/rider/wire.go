//go:build wireinject
// +build wireinject

package rider

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func RiderControllerWire(client *mongo.Database) (*Controller, error) {
	panic(wire.Build(ProviderSet))
}
