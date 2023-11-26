package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/kaytu-io/pennywise/server/cost"
	"github.com/kaytu-io/pennywise/server/resource"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"time"
)

type EchoError struct {
	Message string `json:"message"`
}

type OnboardServiceClient interface {
	GetCost(ctx *echo.Context, sourceID resource.Resource) (*cost.Cost, error)
}

type serverClient struct {
	baseURL string
}

func NewPennywiseServerClient(baseURL string) *serverClient {
	return &serverClient{
		baseURL: baseURL,
	}
}

func (s *serverClient) GetCost(req resource.Resource) (*cost.Cost, error) {
	url := fmt.Sprintf("%s/api/v1/cost/resource", s.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var cost cost.Cost
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
	if res.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(res.Body)
		if err != nil {
			return statusCode, fmt.Errorf("gzip new reader: %w", err)
		}
		defer body.Close()
	}

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
