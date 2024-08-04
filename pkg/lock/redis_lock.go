package lock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisLock struct {
	Client *redis.Client
}

func NewRedisLock(client *redis.Client) *RedisLock {
	return &RedisLock{
		Client: client,
	}
}

func (r *RedisLock) Lock(ctx context.Context, key string, d time.Duration) (bool, error) {
	fmt.Printf("Acquiring lock for " + key)
	res := r.Client.SetNX(ctx, key, "true", d)
	if res.Err() != nil {
		return false, res.Err()
	}
	return res.Val(), nil
}

func (r *RedisLock) Release(ctx context.Context, lockKey string) (bool, error) {
	fmt.Printf("Releasing lock for " + lockKey)
	res := r.Client.Del(ctx, lockKey)
	if res.Err() != nil {
		fmt.Printf("[Cache] failed to release lock. key: %s, err: %s", lockKey, res.Err().Error())
	}
	return res.Val() > 0, res.Err()
}
