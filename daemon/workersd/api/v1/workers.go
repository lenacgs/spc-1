package v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/sibsfps/spc/spc-1/daemon/workersd/api/v1/generated/model"
	"gitlab.com/sibsfps/spc/spc-1/daemon/workersd/lib/context"
	"gitlab.com/sibsfps/spc/spc-1/data"
	"gitlab.com/sibsfps/spc/spc-1/data/transactions"
	"gitlab.com/sibsfps/spc/spc-1/protocol"

	"github.com/labstack/echo/v4"
)

type WorkersInterface interface {
	Process(txn *data.BacklogMsg)
}

// (POST /transaction)
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
		fmt.Sprintf("%v", reqBytes),
	).Info("Debugging bytes")

	dec := protocol.NewDecoder(reqReader)
	var txn transactions.Transaction
	err = dec.Decode(&txn)
	if err != nil {
		log.Errorf("could not decode body: %v", err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	msg := data.BacklogMsg{
		Txn:        txn,
		ReplyQueue: make(chan data.Result),
	}
	defer close(msg.ReplyQueue)

	v1.Node.Process(&msg)

	res := <-msg.ReplyQueue
	if res.Error != nil {
		log.Errorf("error processing transaction: %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	respM := model.Response{}
	for _, m := range res.Mutations {
		item := model.ResponseItem{
			Id:  m.Id,
			Old: m.Old,
			New: m.New,
		}
		respM = append(respM, item)
	}

	var resp bytes.Buffer
	writer := io.Writer(&resp)
	enc := protocol.NewEncoder(writer)
	err = enc.Encode(respM)
	if err != nil {
		log.Errorf("could not encode response: %v", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	log.With(
		context.ResponseBody,
		fmt.Sprintf("%v", resp.Bytes()),
	).Info("Debugging bytes")

	return ctx.Blob(http.StatusOK, "application/x-binary", resp.Bytes())
}
