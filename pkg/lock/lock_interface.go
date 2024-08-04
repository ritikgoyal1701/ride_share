package lock

import (
	"context"
	"time"
)

type Locker interface {
	Lock(ctx context.Context, key string, expiration time.Duration) (bool, error)
	Release(ctx context.Context, lockKey string) (bool, error)
}
