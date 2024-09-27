package tools

import (
	"log"
)

var (
	NodeID int64 = 1
)

func Log(header string, err interface{}) {
	log.Println(header, ":", err)
}

func ConvertStringToInt64(str string) int64 {
	var num int64
	for _, char := range str {
		num = num*10 + int64(char-'0')
	}
	return num
}

func AddIndexToData(arr []byte, index byte) []byte {
	if len(arr) == 0 {
		return []byte{index}
	}

	copy(arr[1:], arr[:len(arr)-1])

	arr[0] = index

	return arr
}
