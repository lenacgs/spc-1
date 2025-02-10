package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	lamport "github.com/ISSuh/logical-clock"
	"net/http"
	"net/url"

	servicemodel "gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1/generated/model"
	workersmodel "gitlab.com/sibsfps/spc/spc-1/daemon/workersd/api/v1/generated/model"
	"gitlab.com/sibsfps/spc/spc-1/protocol"
)

type RestClient struct {
	serverURL  url.URL
	httpClient *http.Client
}

const (
	healthCheckEndpoint = "/v1/health"
	requestEndpoint     = "/v1/request"
	cacheEndpoint       = "/v1/cache"
)

func MakeRestClient(url url.URL) (*RestClient, error) {
	tls := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{}
	tr.TLSClientConfig = tls
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	restClient := &RestClient{
		serverURL:  url,
		httpClient: client,
	}

	resp, err := client.Get(url.String() + healthCheckEndpoint)
	if err != nil {
		return restClient, fmt.Errorf("couldn't connect to the rest client: %s", err)
	}

	if resp.StatusCode != 200 {
		return restClient, fmt.Errorf("couldn't connect to the rest client: status code %d", resp.StatusCode)
	}

	return restClient, nil
}

func encode(v interface{}) (*bytes.Buffer, error) {
	b := []byte{}
	buff := bytes.NewBuffer(b)
	enc := protocol.NewEncoder(buff)
	err := enc.Encode(v)
	if err != nil {
		return nil, fmt.Errorf("could not encode request: %v", err)
	}

	fmt.Printf("encoded request: %v\n", buff.Bytes())

	return buff, nil
}

func (rc *RestClient) doRequest(stmt interface{}) (workersmodel.Response, error) {
	buff, err := encode(stmt)
	if err != nil {
		return nil, err
	}

	urlString := rc.serverURL.String() + requestEndpoint
	resp, err := rc.httpClient.Post(urlString, "application/x-binary", buff)
	if err != nil {
		return nil, fmt.Errorf("could not execute http request: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code response: %v", resp.StatusCode)
	}

	res := workersmodel.Response{}
	dec := protocol.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("could not parse http response: %v", err)
	}

	return res, nil
}

func (rc *RestClient) doRequestCache(stmt interface{}) (servicemodel.Response, error) {
	buff, err := encode(stmt)
	if err != nil {
		return nil, err
	}

	urlString := rc.serverURL.String() + cacheEndpoint
	resp, err := rc.httpClient.Post(urlString, "application/x-binary", buff)
	if err != nil {
		return nil, fmt.Errorf("could not execute http request: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code response: %v", resp.StatusCode)
	}

	res := servicemodel.Response{}
	dec := protocol.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("could not parse http response: %v", err)
	}

	return res, nil
}

func (rc *RestClient) Get(ids []workersmodel.Id) (workersmodel.Response, error) {
	stmt := workersmodel.Select{
		Type: 2,
		Ids:  ids,
	}

	return rc.doRequest(stmt)
}

func (rc *RestClient) Put(ids []workersmodel.Id, statuses []workersmodel.Status) (workersmodel.Response, error) {
	if len(ids) != len(statuses) {
		return nil, fmt.Errorf("missing status")
	}

	workers := make([]workersmodel.Record, 0)
	for i := range ids {
		workers = append(
			workers,
			workersmodel.Record{
				Id:     ids[i],
				Status: statuses[i],
			},
		)
	}

	stmt := workersmodel.Upsert{
		Type:    1,
		Workers: workers,
	}

	return rc.doRequest(stmt)
}

func (rc *RestClient) Del(ids []workersmodel.Id) (workersmodel.Response, error) {
	stmt := workersmodel.Delete{
		Type: 3,
		Ids:  ids,
	}

	return rc.doRequest(stmt)
}

func (rc *RestClient) Cache(ids []servicemodel.Id, clock *lamport.LamportClock) (servicemodel.Response, error) {
	clock.Increase()
	stmt := servicemodel.Request{
		Ids:       ids,
		Timestamp: servicemodel.Timestamp(clock.Time()),
	}

	return rc.doRequestCache(stmt)
}
