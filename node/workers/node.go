package workersNode

import (
	"context"
	"fmt"

	"gitlab.com/sibsfps/spc/spc-1/config"
	"gitlab.com/sibsfps/spc/spc-1/data"
	"gitlab.com/sibsfps/spc/spc-1/logging"

	"github.com/xboshy/go-deadlock"
)

type WorkersNode struct {
	mu         deadlock.Mutex
	ctx        context.Context
	config     config.Local
	cancelCtx  context.CancelFunc
	txnHandler *data.TxnHandler
	log        logging.Logger
}

type StatusReport struct {
}

func MakeNode(log logging.Logger, rootDir string, cfg config.Local) (*WorkersNode, error) {
	var err error

	node := new(WorkersNode)
	node.log = log.With("name", cfg.NetAddress)
	node.config = cfg

	node.txnHandler, err = data.MakeTxnHandler(log)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize the transaction handler: %s", err)
	}

	return node, nil
}

func (node *WorkersNode) Config() config.Local {
	return node.config
}

func (node *WorkersNode) Status() (StatusReport, error) {
	var s StatusReport
	var err error

	return s, err
}

func (node *WorkersNode) Start() {
	node.mu.Lock()
	defer node.mu.Unlock()

	node.ctx, node.cancelCtx = context.WithCancel(context.Background())
	node.txnHandler.Start()
}

func (node *WorkersNode) Stop() {
	node.cancelCtx()
}

func (node *WorkersNode) Process(txn *data.BacklogMsg) {
	node.txnHandler.Process(txn)
}
