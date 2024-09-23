package tools

import (
	"log"
	"os"
)

func Log(t int, header string, err interface{}) {
	switch t {
	default:
	case 0:
		log.Println(header, ":", err)
	case 1:
		log.Fatal(header, ":", err)
		os.Exit(1)
	}
}
