package redis

import "github.com/go-redis/redis/v8"

type Redis struct {
	*redis.Client
}

var redisInstance *Redis

func GetClient() *Redis {
	return redisInstance
}

func SetClient(client *redis.Client) {
	redisInstance = &Redis{client}
}
