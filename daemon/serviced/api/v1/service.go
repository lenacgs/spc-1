package v1

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1/generated/model"
	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/lib/context"
	"gitlab.com/sibsfps/spc/spc-1/data"
	"gitlab.com/sibsfps/spc/spc-1/data/queries"
	"gitlab.com/sibsfps/spc/spc-1/protocol"
	"io"
	"net/http"
)

type ServiceInterface interface {
	Cache(ctx echo.Context) error
	Process(query *data.QueryBacklogMsg)
}

// (POST /cache)
func (v1 *Handlers) RawRequest(ctx echo.Context) error {
	log := v1.Log

	req := ctx.Request()
	if req == nil {
		log.Error("request can't be nil")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	reqBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Errorf("could not read body: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	reqReader := bytes.NewReader(reqBytes)

	log.With(
		context.RequestBody,
		fmt.Sprintf("%+v", reqBytes),
	).Info("debugging bytes")

	dec := protocol.NewDecoder(reqReader)

	var query queries.Query
	err = dec.Decode(&query)
	if err != nil {
		log.Errorf("could not decode body: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	msg := data.QueryBacklogMsg{
		Query:      query,
		ReplyQueue: make(chan data.QueryResult),
	}
	defer close(msg.ReplyQueue)

	v1.Node.Process(&msg)
	res := <-msg.ReplyQueue
	if res.Error != nil {
		log.Errorf("error processing query", res.Error)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	respM := model.Response{}
	for _, m := range res.Statuses {
		item := model.ResponseItem{
			Id:     m.Id,
			Status: m.Location,
		}
		respM = append(respM, item)
	}

	var resp bytes.Buffer
	writer := io.Writer(&resp)
	enc := protocol.NewEncoder(writer)
	err = enc.Encode(&respM)
	if err != nil {
		log.Errorf("could not encode response: %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	log.With(
		context.ResponseBody,
		fmt.Sprintf("%v", resp.Bytes()),
	).Info("debugging bytes")

	return ctx.Blob(http.StatusOK, "application/x-binary", resp.Bytes())

}
