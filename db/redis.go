package db

import (
	"context"
	"sia/backend/tools"

	"github.com/redis/go-redis/v9"
)

type (
	RDB struct {
		*redis.Client
	}
)

var (
	RedisDB *RDB
)

func InitRedis() {
	opt, err := redis.ParseURL(tools.REDIS_URI)
	if err != nil {
		tools.Log("[REDIS]", "Failed to parse Redis URI")
		panic(err)
	}
	client := redis.NewClient(opt)
	RedisDB = &RDB{client}
	tools.Log("[REDIS]", "Connected to Redis")
}

func (h RDB) Get(key string) (string, string) {
	data, err := h.Client.Get(context.Background(), key).Result()
	if err != nil {
		return "", tools.HandleRedisError(err)
	}

	return data, tools.OK
}

func (h RDB) Set(key string, value interface{}) string {
	_, err := h.Client.Set(context.Background(), key, value, 0).Result()
	if err != nil {
		return tools.HandleRedisError(err)
	}

	return tools.OK
}

func (h RDB) Del(key string) string {
	_, err := h.Client.Del(context.Background(), key).Result()
	if err != nil {
		return tools.HandleRedisError(err)
	}

	return tools.OK
}

func (h RDB) Exists(key string) string {
	_, err := h.Client.Exists(context.Background(), key).Result()
	if err != nil {
		return tools.HandleRedisError(err)
	}

	return tools.OK
}

func (h RDB) Close() {
	h.Client.Close()
}
