package services

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/xboshy/go-deadlock"
	"gitlab.com/sibsfps/spc/spc-1/config"
	"gitlab.com/sibsfps/spc/spc-1/data"
	"gitlab.com/sibsfps/spc/spc-1/logging"
)

type ServiceNode struct {
	ctx          context.Context
	cancelCtx    context.CancelFunc
	config       config.Local
	log          logging.Logger
	mu           deadlock.Mutex
	queryHandler *data.QueryHandler
}

type StatusReport struct {
}

func (node *ServiceNode) Start() {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.ctx, node.cancelCtx = context.WithCancel(context.Background())
	node.queryHandler.Start()
}

func (node *ServiceNode) Stop() {
	node.cancelCtx()
}

func (node *ServiceNode) Status() (StatusReport, error) {
	var s StatusReport
	var err error

	return s, err
}

func (node *ServiceNode) Config() config.Local { return node.config }

func MakeNode(log logging.Logger, cfg config.Local) (*ServiceNode, error) {
	var err error

	node := new(ServiceNode)
	node.log = log.With("name", cfg.NetAddress)
	node.config = cfg

	node.queryHandler, err = data.MakeQueryHandler(log, cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize the query handler: %s, er")
	}

	return node, nil
}

func (node *ServiceNode) Process(query *data.QueryBacklogMsg) {
	node.queryHandler.Process(query)
}

func (node *ServiceNode) Cache(ctx echo.Context) error {
	err := node.Cache(ctx)
	if err != nil {
		return err
	}
	return nil
}
