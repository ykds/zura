package snowflake

import "github.com/bwmarrin/snowflake"

var (
	node *snowflake.Node
)

func InitSnowflake(n int64) {
	var err error
	node, err = snowflake.NewNode(n)
	if err != nil {
		panic(err)
	}
}

func NewId() int64 {
	return node.Generate().Int64()
}
