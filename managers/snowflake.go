package managers

import (
	"github.com/bwmarrin/snowflake"
)

// SnowflakeManager represents generating snowflakes.
type SnowflakeManager struct {
	Node *snowflake.Node
}

func NewSnowflakeManager() (*SnowflakeManager, error) {
	// TODO: set this as an env variable
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}

	return &SnowflakeManager{
		Node: node,
	}, nil
}

func (m *SnowflakeManager) Generate() string {
	flake := m.Node.Generate()
	return flake.String()
}
