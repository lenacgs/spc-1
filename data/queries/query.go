package queries

import (
	"gitlab.com/sibsfps/spc/spc-1/protocol"
)

// Query contains the request made by the client (list of ids and timestamp)
type Query struct {
	Timestamp protocol.Timestamp  `codec:"timestamp"`
	Ids       []protocol.WorkerID `codec:"ids"`
}
