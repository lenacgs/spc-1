package clients

import (
	"gitlab.com/sibsfps/spc/spc-1/daemon/workersd/api/v1/generated/model"
)

const workersRequestEndpoint = "/v1/request"

type workers struct {
	restClient restClient
}

func NewWorkers(host string) (*workers, error) {
	restClient, err := makeRestClient(host)
	if err != nil {
		return nil, err

	}

	return &workers{
		restClient: restClient,
	}, nil
}

func (w *workers) Delay(t Time) error {
	return nil
}

func (w *workers) Put(records []Record) ([]Status, error) {
	workers := []model.Record{}
	for _, r := range records {
		workers = append(workers,
			model.Record{
				Id:     r.Id,
				Status: r.Status,
			},
		)
	}

	stmt := model.Upsert{
		Type:    1,
		Workers: workers,
	}

	res, err := w.doRequest(workersRequestEndpoint, stmt)
	if err != nil {
		return nil, err
	}

	mapResponse := map[Id]Status{}
	for _, r := range res {
		mapResponse[Id(r.Id)] = Status(r.New)
	}

	statuses := []Status{}
	for _, r := range records {
		status, ok := mapResponse[Id(r.Id)]
		if !ok {
			status = -1
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (w *workers) Del(ids []Id) ([]Status, error) {
	var mids []model.Id
	for _, id := range ids {
		mids = append(mids, int(id))
	}

	stmt := model.Delete{
		Type: 3,
		Ids:  mids,
	}

	res, err := w.doRequest(workersRequestEndpoint, stmt)
	if err != nil {
		return nil, err
	}

	mapResponse := map[Id]Status{}
	for _, r := range res {
		mapResponse[Id(r.Id)] = Status(r.Old)
	}

	statuses := []Status{}
	for _, id := range ids {
		status, ok := mapResponse[Id(id)]
		if !ok {
			status = -1
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (w *workers) doRequest(endpoint string, stmt interface{}) (model.Response, error) {
	return doRequest[model.Response](w.restClient.httpClient, w.restClient.serverURL, endpoint, stmt)
}
