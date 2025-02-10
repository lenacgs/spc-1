package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1"
	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1/generated/common"
	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1/generated/service"
	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/lib/middlewares"
	"gitlab.com/sibsfps/spc/spc-1/logging"
	node "gitlab.com/sibsfps/spc/spc-1/node/services"
	"net"
)

const (
	BaseURL = "v1"
)

type APINode struct {
	*node.ServiceNode
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
	service.RegisterHandlersWithBaseURL(e, &v1Handler, BaseURL)

	return e
}
