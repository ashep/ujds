package client

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"

	ujdsproto "github.com/ashep/ujds/sdk/proto/ujds/v1"
	ujdsconnect "github.com/ashep/ujds/sdk/proto/ujds/v1/v1connect"
)

type Client struct {
	I ujdsconnect.IndexServiceClient
	R ujdsconnect.RecordServiceClient
}

func New(url, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	icp := connect.WithInterceptors(NewAuthInterceptor(apiKey))

	return &Client{
		I: ujdsconnect.NewIndexServiceClient(httpClient, url, icp),
		R: ujdsconnect.NewRecordServiceClient(httpClient, url, icp),
	}
}

func (c *Client) IndexExists(ctx context.Context, name string) (bool, error) {
	_, err := c.I.GetIndex(ctx, connect.NewRequest(&ujdsproto.GetIndexRequest{Name: name}))
	if err != nil && connect.CodeOf(err) == connect.CodeNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
