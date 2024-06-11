package transactions

import (
	"gitlab.com/sibsfps/spc/spc-1/data/workers"
	"gitlab.com/sibsfps/spc/spc-1/protocol"
)

type Transaction struct {
	Type    protocol.TxnType    `codec:"type"`
	Ids     []protocol.WorkerID `codec:"ids"`
	Workers []workers.Worker    `codec:"workers"`
}

type WorkerMutation struct {
	Id  protocol.WorkerID `codec:"id"`
	Old protocol.Status   `codec:"old"`
	New protocol.Status   `codec:"new"`
}
