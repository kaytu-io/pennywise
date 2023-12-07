package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
	"time"
)

type EchoError struct {
	Message string `json:"message"`
}

type OnboardServiceClient interface {
	GetResourceCost(req resource.Resource) (*cost.Cost, error)
	GetStateCost(req []resource.Resource) (*cost.Cost, error)
	IngestAws(service, region string) error
	IngestAzure(service, region string) error
}

type serverClient struct {
	baseURL string
}

func NewPennywiseServerClient(baseURL string) *serverClient {
	return &serverClient{
		baseURL: baseURL,
	}
}

func (s *serverClient) IngestAws(service, region string) error {
	url := fmt.Sprintf("%s/api/v1/ingest/aws?service=%s&region=%s", s.baseURL, service, region)
	url = strings.ReplaceAll(url, " ", "%20")
	if statusCode, err := doRequest(http.MethodPut, url, nil, nil); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}
	return nil
}

func (s *serverClient) IngestAzure(service, region string) error {
	url := fmt.Sprintf("%s/api/v1/ingest/azure?service=%s&region=%s", s.baseURL, service, region)
	url = strings.ReplaceAll(url, " ", "%20")
	if statusCode, err := doRequest(http.MethodPut, url, nil, nil); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}
	return nil
}

func (s *serverClient) GetResourceCost(req resource.Resource) (*cost.State, error) {
	url := fmt.Sprintf("%s/api/v1/cost/resource", s.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var cost cost.State
	if statusCode, err := doRequest(http.MethodGet, url, payload, &cost); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &cost, nil
}

func (s *serverClient) GetStateCost(req resource.State) (*cost.State, error) {
	url := fmt.Sprintf("%s/api/v1/cost/state", s.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var cost cost.State
	if statusCode, err := doRequest(http.MethodGet, url, payload, &cost); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &cost, nil
}

func doRequest(method, url string, payload []byte, v interface{}) (statusCode int, err error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return statusCode, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set(echo.HeaderContentType, "application/json")
	t := http.DefaultTransport.(*http.Transport)
	client := http.Client{
		Timeout:   3 * time.Minute,
		Transport: t,
	}

	res, err := client.Do(req)
	if err != nil {
		return statusCode, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()
	body := res.Body

	statusCode = res.StatusCode
	if res.StatusCode != http.StatusOK {
		d, err := io.ReadAll(body)
		if err != nil {
			return statusCode, fmt.Errorf("read body: %w", err)
		}

		var echoerr EchoError
		if jserr := json.Unmarshal(d, &echoerr); jserr == nil {
			return statusCode, fmt.Errorf(echoerr.Message)
		}

		return statusCode, fmt.Errorf("http status: %d: %s", res.StatusCode, d)
	}
	if v == nil {
		return statusCode, nil
	}

	return statusCode, json.NewDecoder(body).Decode(v)
}
