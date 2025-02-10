package data

import (
	"context"
	"gitlab.com/sibsfps/spc/spc-1/config"
	"gitlab.com/sibsfps/spc/spc-1/data/queries"
	"gitlab.com/sibsfps/spc/spc-1/logging"
	"sync"
)

type QueryHandler struct {
	ctx          context.Context
	ctxCancel    context.CancelFunc
	backlogWg    sync.WaitGroup
	backlogQueue chan *QueryBacklogMsg
	cache        Cache
	log          logging.Logger
}

type QueryBacklogMsg struct {
	Query      queries.Query
	ReplyQueue chan QueryResult
}

type QueryResult struct {
	Statuses []cacheItem
	Error    error
}

func MakeQueryHandler(log logging.Logger, cfg config.Local) (*QueryHandler, error) {
	backlogSize := 10

	cache, err := MakeCache(log, cfg)
	if err != nil {
		return nil, err
	}

	handler := &QueryHandler{
		cache:        cache,
		backlogQueue: make(chan *QueryBacklogMsg, backlogSize),
		log:          log,
	}

	return handler, nil
}

func (handler *QueryHandler) Start() {
	handler.ctx, handler.ctxCancel = context.WithCancel(context.Background())
	handler.backlogWg.Add(1)
	go handler.handler()
}

func (handler *QueryHandler) handler() {
	defer handler.backlogWg.Done()
	for {
		select {
		case msg := <-handler.backlogQueue:
			reply := QueryResult{}

			// statuses is an array of cacheItem{Location, Expiration, Id}
			statuses, err := handler.cache.Query(msg.Query)

			reply.Statuses = statuses
			reply.Error = err

			msg.ReplyQueue <- reply
		case <-handler.ctx.Done():
			return
		}
	}
}

func (handler *QueryHandler) Process(query *QueryBacklogMsg) {
	handler.backlogQueue <- query
}

func (handler *QueryHandler) Cache() Cache {
	return handler.cache
}
