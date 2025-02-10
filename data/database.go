package data

import (
	"gitlab.com/sibsfps/spc/spc-1/data/transactions"
	"gitlab.com/sibsfps/spc/spc-1/data/workers"
	"gitlab.com/sibsfps/spc/spc-1/logging"
	"gitlab.com/sibsfps/spc/spc-1/protocol"

	"github.com/xboshy/go-deadlock"
)

type Database interface {
	Upsert([]workers.Worker) ([]transactions.WorkerMutation, error)
	Select([]protocol.WorkerID) ([]transactions.WorkerMutation, error)
	Delete([]protocol.WorkerID) ([]transactions.WorkerMutation, error)
}

type database struct {
	mu  deadlock.Mutex
	log logging.Logger

	data map[protocol.WorkerID]protocol.Location
}

func MakeDatabase(log logging.Logger) (Database, error) {
	db := new(database)
	db.log = log.With("database", "internal")
	db.data = make(map[protocol.WorkerID]protocol.Location)

	return db, nil
}

func (node *database) Upsert(ws []workers.Worker) ([]transactions.WorkerMutation, error) {
	node.mu.Lock()
	defer node.mu.Unlock()

	res := make([]transactions.WorkerMutation, 0)
	var ok bool
	for _, new := range ws {
		old := 0
		old, ok = node.data[new.Id]
		node.data[new.Id] = new.Status

		if !ok {
			old = protocol.UnavailableStatus
		}

		mutation := transactions.WorkerMutation{
			Id:  new.Id,
			Old: old,
			New: new.Status,
		}

		res = append(res, mutation)
	}

	return res, nil
}

func (node *database) Select(idx []protocol.WorkerID) ([]transactions.WorkerMutation, error) {
	node.mu.Lock()
	defer node.mu.Unlock()

	res := make([]transactions.WorkerMutation, 0)
	var ok bool
	for _, id := range idx {
		cur := 0
		cur, ok = node.data[id]

		if !ok {
			cur = protocol.UnavailableStatus
		}

		mutation := transactions.WorkerMutation{
			Id:  id,
			Old: cur,
			New: cur,
		}

		res = append(res, mutation)
	}

	return res, nil
}

func (node *database) Delete(idx []protocol.WorkerID) ([]transactions.WorkerMutation, error) {
	node.mu.Lock()
	defer node.mu.Unlock()

	res := make([]transactions.WorkerMutation, 0)
	var ok bool
	for _, id := range idx {
		old := 0
		old, ok = node.data[id]
		delete(node.data, id)

		if !ok {
			old = protocol.UnavailableStatus
		}

		mutation := transactions.WorkerMutation{
			Id:  id,
			Old: old,
			New: 0,
		}

		res = append(res, mutation)
	}

	return res, nil
}
