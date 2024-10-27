package lib

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	// Environment variables
	UDP_ADDR string
	UDP_PORT int64
)

func InitEnvVars() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	UDP_ADDR = LoadEnvVar("UDP_ADDR")
	UDP_PORT, err = strconv.ParseInt(LoadEnvVar("UDP_PORT"), 10, 0)
	if err != nil {
		panic(err)
	}
}

func LoadEnvVar(key string) string {
	str := os.Getenv(key)
	if str == "" {
		panic("Environment variable not found: " + key)
	}
	return str
}
