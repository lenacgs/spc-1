package workers

import "gitlab.com/sibsfps/spc/spc-1/protocol"

type Worker struct {
	Id     protocol.WorkerID `codec:"id"`
	Status protocol.Status   `codec:"status"`
}
