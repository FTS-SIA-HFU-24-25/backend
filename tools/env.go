package tools

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	// UDP Server Variables
	UDP_PASSWORD string
	UDP_SALT     string
	PORT         string
	UDP_HOST     string
	UDP_ADDR     string

	// Redis Variables
	REDIS_HOST     string
	REDIS_PASSWORD string

	// UUID Variables
	UUID_HASH_SALT string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	UDP_PASSWORD = os.Getenv("UDP_PASSWORD")
	UDP_SALT = os.Getenv("UDP_SALT")
	PORT = os.Getenv("UDP_PORT")
	UDP_HOST = os.Getenv("UDP_HOST")
	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	UUID_HASH_SALT = os.Getenv("UUID_HASH_SALT")
	Log("[ENV]", UDP_HOST+":"+PORT)
	UDP_ADDR = UDP_HOST + ":" + PORT
	Log("[ENV]", "Loaded environment variables...")
}
