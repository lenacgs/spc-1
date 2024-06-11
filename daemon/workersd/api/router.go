package api

import (
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	v1 "gitlab.com/sibsfps/spc/spc-1/daemon/workersd/api/v1"
	"gitlab.com/sibsfps/spc/spc-1/daemon/workersd/api/v1/generated/common"
	"gitlab.com/sibsfps/spc/spc-1/daemon/workersd/api/v1/generated/workers"
	"gitlab.com/sibsfps/spc/spc-1/daemon/workersd/lib/middlewares"
	"gitlab.com/sibsfps/spc/spc-1/logging"
	node "gitlab.com/sibsfps/spc/spc-1/node/workers"
)

const (
	BaseURL = "v1"
)

type APINode struct {
	*node.WorkersNode
}

type APINodeInterface interface {
	v1.NodeInterface
}

func NewRouter(logger logging.Logger, node APINodeInterface, shutdown <-chan struct{}, listener net.Listener) *echo.Echo {
	e := echo.New()
	e.Logger = logger.MakeEchoLogger()

	e.Listener = listener
	e.HideBanner = true

	e.Pre(
		middleware.RemoveTrailingSlash(),
	)
	e.Use(
		middlewares.MakeLogger(logger),
	)

	v1Handler := v1.Handlers{
		Node:     node,
		Log:      logger,
		Shutdown: shutdown,
	}

	common.RegisterHandlersWithBaseURL(e, &v1Handler, BaseURL)
	workers.RegisterHandlersWithBaseURL(e, &v1Handler, BaseURL)

	return e
}
