package worker

import (
	"context"
	"errors"

	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

type disabledHTTPClient struct{}

func (disabledHTTPClient) Do(context.Context, *ferrethttp.Request) (*ferrethttp.Response, error) {
	return nil, errors.New("http: outbound requests are disabled")
}
