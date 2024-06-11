package v1

import "gitlab.com/sibsfps/spc/spc-1/logging"

// NodeInterface represents node fns used by the handlers.
type NodeInterface interface {
	CommonInterface
	WorkersInterface
}

type Handlers struct {
	Node     NodeInterface
	Log      logging.Logger
	Shutdown <-chan struct{}
}
