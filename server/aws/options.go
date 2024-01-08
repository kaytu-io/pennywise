package aws

import (
	"net/http"
)

//go:generate mockgen -destination=../mock/http_client.go -mock_names=HTTPClient=HTTPClient -package mock github.com/kaytu-io/pennywise/server/aws HTTPClient

// HTTPClient is an interface of a client that is able to Do HTTP requests
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Option is used to configure the Ingester.
type Option func(ing *Ingester)
