package tools

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	OK                = "OK"
	REDIS_NOT_FOUND   = "REDIS_NOT_FOUND"
	REDIS_GET_FAILED  = "REDIS_GET_FAILED"
	REDIS_UNKNOWN_ERR = "REDIS_UNKNOWN_ERR"
	DB_DUP_KEY        = "DB_DUP_KEY"
	DB_REC_NOTFOUND   = "DB_REC_NOTFOUND"
	DB_UNKNOWN_ERR    = "DB_UNKNOWN_ERR"
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

func DBHandleError(e error) string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		Log("[DATABASE]", fmt.Sprintf("Error from %s, line %d", file, line))
	}
	var res string
	if errors.Is(e, gorm.ErrDuplicatedKey) {
		res = DB_DUP_KEY
	} else if errors.Is(e, gorm.ErrRecordNotFound) {
		res = DB_REC_NOTFOUND
	} else {
		res = DB_UNKNOWN_ERR
	}
	Log("[DATABASE]", e.Error())
	return res
}
