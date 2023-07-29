package astutils

import "fmt"

type (
	NodeId  uint64
	Version uint64

	NodeIded interface {
		ID() NodeId
	}

	Versioned interface {
		Version() Version
	}

	NodeIdManager struct {
		nodeId NodeId
	}

	VersionManager struct {
		version Version
	}
)

func NewNodeIdManager() NodeIdManager {
	return NodeIdManager{0}
}

func (m *NodeIdManager) Next() NodeId {
	m.nodeId++
	return m.nodeId
}

func NewVersionManager() VersionManager {
	return VersionManager{0}
}

func (m *VersionManager) Next() Version {
	m.version++
	return m.version
}

func (n NodeId) String() string {
	return fmt.Sprintf("%d", n)
}

func (v Version) String() string {
	return fmt.Sprintf("%d", v)
}
