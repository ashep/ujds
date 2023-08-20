package client

import (
	"net/http"

	"github.com/bufbuild/connect-go"

	indexconnect "github.com/ashep/ujds/sdk/proto/ujds/index/v1/v1connect"
	recordconnect "github.com/ashep/ujds/sdk/proto/ujds/record/v1/v1connect"
)

type Client struct {
	I indexconnect.IndexServiceClient
	R recordconnect.RecordServiceClient
}

func New(url, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	icp := connect.WithInterceptors(NewAuthInterceptor(apiKey))

	return &Client{
		I: indexconnect.NewIndexServiceClient(httpClient, url, icp),
		R: recordconnect.NewRecordServiceClient(httpClient, url, icp),
	}
}
