package v1

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"gitlab.com/sibsfps/spc/spc-1/config"
	node "gitlab.com/sibsfps/spc/spc-1/node/services"
	"net/http"
)

type CommonInterface interface {
	Status() (s node.StatusReport, err error)
	Config() config.Local
}

func (h *Handlers) HealthCheck(ctx echo.Context) error {
	w := ctx.Response().Writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct{}{})

	return nil
}
