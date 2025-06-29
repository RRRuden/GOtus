package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLogger struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisLogger(addr string, password string, db int) Logger {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	return &RedisLogger{
		client: rdb,
		ctx:    ctx,
	}
}

func (r *RedisLogger) Log(key string, message string, ttlSeconds int) error {
	err := r.client.Set(r.ctx, key, message, time.Duration(ttlSeconds)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("ошибка логгирования в Redis: %w", err)
	}
	return nil
}
