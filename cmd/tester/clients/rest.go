package clients

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	smodel "gitlab.com/sibsfps/spc/spc-1/daemon/serviced/api/v1/generated/model"
	wmodel "gitlab.com/sibsfps/spc/spc-1/daemon/workersd/api/v1/generated/model"
	"gitlab.com/sibsfps/spc/spc-1/protocol"
)

const healthCheckEndpoint = "/v1/health"

type restClient struct {
	serverURL  string
	httpClient *http.Client
}

func makeRestClient(urlString string) (restClient, error) {
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = fmt.Sprintf("http://%s", urlString)
	}

	url, err := url.Parse(urlString)
	if err != nil {
		log.Println("Warning couldn't parse URL", strconv.Quote(err.Error()))
	}

	tls := &tls.Config{
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{}
	tr.TLSClientConfig = tls
	httpClient := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	restClient := restClient{
		serverURL:  url.String(),
		httpClient: httpClient,
	}

	resp, err := httpClient.Get(restClient.serverURL + healthCheckEndpoint)
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

	return buff, nil
}

func doRequest[T smodel.Response | wmodel.Response](httpClient *http.Client, baseurl string, endpoint string, stmt interface{}) (T, error) {
	buff, err := encode(stmt)
	if err != nil {
		return nil, err
	}

	url := baseurl + endpoint
	resp, err := httpClient.Post(url, "application/x-binary", buff)
	if err != nil {
		return nil, fmt.Errorf("could not execute http request: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code response: %v", resp.StatusCode)
	}

	var res T
	dec := protocol.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("could not parse http response: %v", err)
	}

	return res, nil
}
