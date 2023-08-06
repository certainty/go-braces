package astutils

import "fmt"

type (
	NodeId  string
	Version uint64

	NodeIded interface {
		ID() NodeId
	}

	Versioned interface {
		Version() Version
	}

	NodeIdManager struct {
		counter uint64
		prefix  string
	}

	VersionManager struct {
		version Version
	}
)

func NewNodeIdManager(prefix string) NodeIdManager {
	return NodeIdManager{prefix: prefix}
}

func (m *NodeIdManager) Next() NodeId {
	m.counter++
	return NodeId(fmt.Sprintf("%s%d", m.prefix, m.counter))
}

func NewVersionManager() VersionManager {
	return VersionManager{0}
}

func (m *VersionManager) Next() Version {
	m.version++
	return m.version
}

func (v Version) String() string {
	return fmt.Sprintf("%d", v)
}
