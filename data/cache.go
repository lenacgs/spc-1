package data

import (
	"fmt"
	lamport "github.com/ISSuh/logical-clock"
	"github.com/xboshy/go-deadlock"
	"gitlab.com/sibsfps/spc/spc-1/config"
	"gitlab.com/sibsfps/spc/spc-1/daemon/workersd/api/v1/generated/model"
	"gitlab.com/sibsfps/spc/spc-1/data/queries"
	"gitlab.com/sibsfps/spc/spc-1/data/requests"
	"gitlab.com/sibsfps/spc/spc-1/logging"
	"gitlab.com/sibsfps/spc/spc-1/protocol"
)

type Cache interface {
	Query(query queries.Query) ([]cacheItem, error)
	apiRequest([]protocol.WorkerID) (model.Response, error)
	calculateExpiration(cacheItem) protocol.Timestamp
	add(cacheItem) cacheItem
	evict()
	get(protocol.WorkerID) (cacheItem, bool)
}

const (
	workersApiUrl = "http://localhost:8080"
)

type cache struct {
	mu  deadlock.Mutex
	log logging.Logger

	// cache data
	data map[protocol.WorkerID]cacheItem

	clock *lamport.LamportClock

	unavailableTTL uint64
	locationTTL    uint64
	maxCapacity    int
	size           int

	/* saves the order in which workers are added to the cache
	(for FIFO deletion when cache reaches max capacity) */
	creationOrder []protocol.WorkerID

	restClient *requests.RestClient
}

type cacheItem struct {
	Location   protocol.Location
	Expiration protocol.Timestamp
	Id         protocol.WorkerID
}

// MakeCache makes new cache according to config
func MakeCache(log logging.Logger, cfg config.Local) (Cache, error) {
	var err error

	cache := new(cache)
	cache.log = log.With("cache", "internal")
	cache.data = make(map[protocol.WorkerID]cacheItem)
	cache.clock = lamport.NewLamportClock()
	cache.size = 0
	cache.maxCapacity = cfg.CacheMaxCapacity
	cache.unavailableTTL = uint64(cfg.CacheUnavailableTTL)
	cache.locationTTL = uint64(cfg.CacheLocationTTL)
	cache.creationOrder = make([]protocol.WorkerID, cache.maxCapacity)

	cache.restClient, err = requests.MakeRestClient(*requests.MakeUrl(workersApiUrl))
	if err != nil {
		cache.log.Errorf("error making rest client", err.Error())
	}

	return cache, nil
}

// calculateExpiration calculates the tick in which a worker's location becomes invalid
func (cache *cache) calculateExpiration(worker cacheItem) protocol.Timestamp {
	switch worker.Location {
	case 0:
		return cache.clock.Time() + cache.unavailableTTL
	default:
		return cache.clock.Time() + cache.locationTTL
	}
}

// add adds a new worker's location to the cache, taking into account the cache max capacity
func (cache *cache) add(worker cacheItem) cacheItem {
	cache.log.Infof("add: ", worker.Id)
	cache.mu.Lock()
	defer cache.mu.Unlock()

	expiration := cache.calculateExpiration(worker)

	worker.Expiration = expiration
	cache.log.Infof("expiration for this entry:", expiration)

	// if cache has reached max capacity, delete oldest record
	if cache.size+1 > cache.maxCapacity {
		cache.evict()
	}

	cache.data[worker.Id] = worker
	cache.size++
	cache.creationOrder = append(cache.creationOrder, worker.Id)

	return worker
}

// evict deletes the oldest cache item
func (cache *cache) evict() {
	cache.log.Infof("evict: ", cache.creationOrder)
	deletable := cache.creationOrder[0]
	cache.creationOrder = cache.creationOrder[1:]
	delete(cache.data, deletable)
	cache.size--
}

// get gets worker location from the cache, checking its expiration
func (cache *cache) get(worker protocol.WorkerID) (cacheItem, bool) {
	cache.log.Infof("get: ", worker)
	cache.mu.Lock()
	defer cache.mu.Unlock()

	item, ok := cache.data[worker]
	cache.log.Infof("getting cached value: ", item, ok)
	if item.Expiration <= cache.clock.Time() {
		return item, false
	}
	return item, ok
}

// Query receives the client's request for workers location and returns the corresponding cache records
func (cache *cache) Query(query queries.Query) ([]cacheItem, error) {
	cache.log.Infof("query: ", query)

	// update service clock with the clock from the incoming query
	cache.clock.Update(query.Timestamp)
	cache.log.Info("service clock = ", cache.clock)

	items := make([]cacheItem, 0)
	toFetch := make([]protocol.WorkerID, 0)
	var item cacheItem
	var ok bool

	for _, workerId := range query.Ids {
		item, ok = cache.get(workerId)

		// if cache miss, fetch from worker API
		if !ok {
			cache.log.Infof("cache miss for worker ", workerId)
			toFetch = append(toFetch, workerId)
		} else {
			cache.log.Infof("cache hit for worker ", workerId)
			items = append(items, item)
		}
	}

	// fetching cache misses from workers API
	if len(toFetch) > 0 {
		fetched, err := cache.apiRequest(toFetch)

		// if not possible to fetch, then workers' status is unavailable
		if err != nil {
			for _, workerId := range toFetch {
				item := cache.add(cacheItem{Id: workerId, Location: protocol.UnavailableStatus})
				items = append(items, item)
			}
		} else {
			// adding newly fetched to cache
			for _, responseItem := range fetched {
				item := cache.add(cacheItem{Id: responseItem.Id, Location: responseItem.New})
				items = append(items, item)
			}
		}
	}

	return items, nil
}

// apiRequest makes a request for workers locations to the Workers API
func (cache *cache) apiRequest(workerIDs []protocol.WorkerID) (model.Response, error) {

	cache.log.Infof("api request: ", workerIDs)
	res := make(model.Response, 0)

	if !requests.CheckConnected(cache.restClient) {
		cache.log.Errorf("could not connect to workers API")
		return res, fmt.Errorf("could not connect to workers API")
	}

	res, err := cache.restClient.Get(workerIDs)
	if err != nil {
		cache.log.Errorf("error making http request", err.Error())
		return res, err
	}

	return res, nil
}
