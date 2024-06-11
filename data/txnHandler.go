package data

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/sibsfps/spc/spc-1/data/transactions"
	"gitlab.com/sibsfps/spc/spc-1/logging"
	"gitlab.com/sibsfps/spc/spc-1/protocol"
)

type TxnHandler struct {
	ctx          context.Context
	ctxCancel    context.CancelFunc
	backlogWg    sync.WaitGroup
	backlogQueue chan *BacklogMsg
	db           Database
	log          logging.Logger
}

type BacklogMsg struct {
	Txn        transactions.Transaction
	ReplyQueue chan Result
}

type Result struct {
	Mutations []transactions.WorkerMutation
	Error     error
}

func MakeTxnHandler(log logging.Logger) (*TxnHandler, error) {
	backlogSize := 10

	db, err := MakeDatabase(log)
	if err != nil {
		return nil, err
	}

	handler := &TxnHandler{
		db:           db,
		backlogQueue: make(chan *BacklogMsg, backlogSize),
		log:          log,
	}

	return handler, nil
}

func (handler *TxnHandler) Start() {
	handler.ctx, handler.ctxCancel = context.WithCancel(context.Background())
	handler.backlogWg.Add(1)
	go handler.handler()
}

func (handler *TxnHandler) Stop() {
	handler.ctxCancel()
	handler.backlogWg.Wait()
}

func (handler *TxnHandler) Process(txn *BacklogMsg) {
	handler.backlogQueue <- txn
}

func (handler *TxnHandler) handler() {
	defer handler.backlogWg.Done()
	for {
		var err error

		select {
		case msg := <-handler.backlogQueue:
			var mutations []transactions.WorkerMutation
			reply := Result{}

			switch msg.Txn.Type {
			case protocol.UpsertType:
				mutations, err = handler.db.Upsert(msg.Txn.Workers)
			case protocol.SelectType:
				mutations, err = handler.db.Select(msg.Txn.Ids)
			case protocol.DeleteType:
				mutations, err = handler.db.Delete(msg.Txn.Ids)
			default:
				err = fmt.Errorf("invalid transaction type")
			}

			reply.Mutations = mutations
			reply.Error = err

			msg.ReplyQueue <- reply
		case <-handler.ctx.Done():
			return
		}
	}
}
