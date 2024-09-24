package tools

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

var (
	NodeID int64 = 1
)

func Log(header string, err interface{}) {
	log.Println(header, ":", err)
}

func GenerateUUID() [8]byte {
	node, err := snowflake.NewNode(NodeID)
	if err != nil {
		return [8]byte{}
	}

	id := node.Generate()
	NodeID++
	return id.IntBytes()
}

func AddIndexToData(arr []byte, index byte) []byte {
	if len(arr) == 0 {
		return []byte{index}
	}

	copy(arr[1:], arr[:len(arr)-1])

	arr[0] = index

	return arr
}
