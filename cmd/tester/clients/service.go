package clients

import (
	"gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1/generated/model"
)

const serviceRequestEndpoint = "/v1/cache"

type service struct {
	timestamp  Time
	restClient restClient
}

func NewService(host string) (*service, error) {
	restClient, err := makeRestClient(host)
	if err != nil {
		return nil, err
	}

	return &service{
		timestamp:  0,
		restClient: restClient,
	}, nil
}

func (s *service) Get(ids []Id) ([]Status, error) {
	var mids []model.Id
	for _, id := range ids {
		mids = append(mids, int(id))
	}

	stmt := model.Request{
		Timestamp: int(s.timestamp),
		Ids:       mids,
	}

	res, err := s.doRequest(serviceRequestEndpoint, stmt)
	if err != nil {
		return nil, err
	}

	mapResponse := map[Id]Status{}
	for _, r := range res {
		mapResponse[Id(r.Id)] = Status(r.Status)
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

func (s *service) Forward(t Time) error {
	s.timestamp += t
	return nil
}

func (s *service) doRequest(endpoint string, stmt interface{}) (model.Response, error) {
	return doRequest[model.Response](s.restClient.httpClient, s.restClient.serverURL, endpoint, stmt)
}
