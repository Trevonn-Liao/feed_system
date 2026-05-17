package snowflake

import (
	"sync"

	bwsnowflake "github.com/bwmarrin/snowflake"
)

type Generator struct {
	node *bwsnowflake.Node
	mu   sync.Mutex
}

func New(workerID, datacenterID int64) (*Generator, error) {
	_ = datacenterID
	node, err := bwsnowflake.NewNode(workerID)
	if err != nil {
		return nil, err
	}
	return &Generator{node: node}, nil
}

func (g *Generator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return int64(g.node.Generate())
}
