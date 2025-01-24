package client

import (
	"net/http"

	"connectrpc.com/connect"

	indexconnect "github.com/ashep/ujds/sdk/proto/ujds/index/v1/v1connect"
	recordconnect "github.com/ashep/ujds/sdk/proto/ujds/record/v1/v1connect"
)

type Client struct {
	I indexconnect.IndexServiceClient
	R recordconnect.RecordServiceClient
}

func New(url, authToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	icp := connect.WithInterceptors(NewAuthInterceptor(authToken))

	return &Client{
		I: indexconnect.NewIndexServiceClient(httpClient, url, icp),
		R: recordconnect.NewRecordServiceClient(httpClient, url, icp),
	}
}
