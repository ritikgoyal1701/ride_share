//go:build wireinject
// +build wireinject

package ride

import (
	redis2 "github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func RideControllerWire(client *mongo.Database, redisClient *redis2.Client) (*Controller, error) {
	panic(wire.Build(ProviderSet))
}
