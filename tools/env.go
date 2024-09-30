package tools

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	// Server Variables
	SERVICE_VERSION int

	// UDP Server Variables
	UDP_PASSWORD string
	UDP_SALT     string
	UDP_PORT     string
	UDP_HOST     string
	UDP_ADDR     string

	// TCP Server Variables
	TCP_PORT string

	// Redis Variables
	REDIS_HOST     string
	REDIS_PASSWORD string
	REDIS_PORT     string
	REDIS_ADDR     string
	REDIS_URI      string

	// PostgreSQL Database
	POSTGRES_URI string

	// UUID Variables
	UUID_HASH_SALT string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	SERVICE_VERSION, err = strconv.Atoi(os.Getenv("SERVICE_VERSION"))
	if err != nil {
		panic(err)
	}
	UDP_PASSWORD = os.Getenv("UDP_PASSWORD")
	UDP_SALT = os.Getenv("UDP_SALT")
	UDP_PORT = os.Getenv("UDP_PORT")
	UDP_HOST = os.Getenv("UDP_HOST")
	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	REDIS_PORT = os.Getenv("REDIS_PORT")
	REDIS_URI = os.Getenv("REDIS_URI")
	POSTGRES_URI = os.Getenv("POSTGRES_URI")
	UUID_HASH_SALT = os.Getenv("UUID_HASH_SALT")
	UDP_ADDR = UDP_HOST + ":" + UDP_PORT
	REDIS_ADDR = REDIS_HOST + ":" + REDIS_PORT
	Log("[ENV]", "Loaded environment variables...")
}
