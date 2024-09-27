package uuid

import (
	"github.com/bwmarrin/snowflake"
)

func GenerateUUID(serialID int64) [8]byte {
	node, err := snowflake.NewNode(serialID)
	if err != nil {
		return [8]byte{}
	}

	id := node.Generate()
	return id.IntBytes()
}
