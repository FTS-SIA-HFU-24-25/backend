package tools

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	UDP_PASSWORD string
	UDP_SALT     string
	PORT         string
	UDP_HOST     string
	UDP_ADDR     string
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	UDP_PASSWORD = os.Getenv("UDP_PASSWORD")
	UDP_SALT = os.Getenv("UDP_SALT")
	PORT = os.Getenv("PORT")
	UDP_HOST = os.Getenv("UDP_HOST")
	UDP_ADDR = UDP_HOST + ":" + PORT
}
