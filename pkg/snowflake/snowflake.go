package snowflake

import "github.com/bwmarrin/snowflake"

var (
	node *snowflake.Node
)

func InitSnowflake() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

func NewId() int64 {
	return node.Generate().Int64()
}
