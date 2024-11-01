package lib

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	// Environment variables
	UDP_ADDR string
	UDP_PORT int

	// TCP Vars
	TCP_ADDR string
	TCP_PORT int

	// Websocket
	WEBSOCKET_PATH string
)

func InitEnvVars() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	UDP_ADDR = LoadEnvVar("UDP_ADDR")
	UDP_PORT, err = strconv.Atoi(LoadEnvVar("UDP_PORT"))
	TCP_ADDR = LoadEnvVar("TCP_ADDR")
	TCP_PORT, err = strconv.Atoi(LoadEnvVar("TCP_PORT"))
	WEBSOCKET_PATH = LoadEnvVar("WEBSOCKET_PATH")
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
