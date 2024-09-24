package tools

import "github.com/redis/go-redis/v9"

var (
	OK                = "OK"
	REDIS_NOT_FOUND   = "REDIS_NOT_FOUND"
	REDIS_GET_FAILED  = "REDIS_GET_FAILED"
	REDIS_UNKNOWN_ERR = "REDIS_UNKNOWN_ERR"
)

func HandleRedisError(err error) string {
	if err != nil {
		Log("[REDIS]", err)
		return REDIS_GET_FAILED
	} else if err == redis.Nil {
		Log("[REDIS]", "Key not found")
		return REDIS_NOT_FOUND
	} else {
		return REDIS_UNKNOWN_ERR
	}
}
