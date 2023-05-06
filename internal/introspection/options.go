package introspection

import (
	"context"
	"net"
)

type Options struct {
	GrpServerAddress net.Addr
	API              API
	Cancel           context.CancelFunc
}
