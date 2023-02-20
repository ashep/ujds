// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/datapimp/v1/service.proto

package datapimpv1connect

import (
	context "context"
	errors "errors"
	v1 "github.com/ashep/datapimp/gen/proto/datapimp/v1"
	connect_go "github.com/bufbuild/connect-go"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// ItemServiceName is the fully-qualified name of the ItemService service.
	ItemServiceName = "datapimp.v1.ItemService"
)

// ItemServiceClient is a client for the datapimp.v1.ItemService service.
type ItemServiceClient interface {
	PushItem(context.Context, *connect_go.Request[v1.PushItemRequest]) (*connect_go.Response[v1.PushItemResponse], error)
	GetItem(context.Context, *connect_go.Request[v1.GetItemRequest]) (*connect_go.Response[v1.GetItemResponse], error)
}

// NewItemServiceClient constructs a client for the datapimp.v1.ItemService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewItemServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ItemServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &itemServiceClient{
		pushItem: connect_go.NewClient[v1.PushItemRequest, v1.PushItemResponse](
			httpClient,
			baseURL+"/datapimp.v1.ItemService/PushItem",
			opts...,
		),
		getItem: connect_go.NewClient[v1.GetItemRequest, v1.GetItemResponse](
			httpClient,
			baseURL+"/datapimp.v1.ItemService/GetItem",
			opts...,
		),
	}
}

// itemServiceClient implements ItemServiceClient.
type itemServiceClient struct {
	pushItem *connect_go.Client[v1.PushItemRequest, v1.PushItemResponse]
	getItem  *connect_go.Client[v1.GetItemRequest, v1.GetItemResponse]
}

// PushItem calls datapimp.v1.ItemService.PushItem.
func (c *itemServiceClient) PushItem(ctx context.Context, req *connect_go.Request[v1.PushItemRequest]) (*connect_go.Response[v1.PushItemResponse], error) {
	return c.pushItem.CallUnary(ctx, req)
}

// GetItem calls datapimp.v1.ItemService.GetItem.
func (c *itemServiceClient) GetItem(ctx context.Context, req *connect_go.Request[v1.GetItemRequest]) (*connect_go.Response[v1.GetItemResponse], error) {
	return c.getItem.CallUnary(ctx, req)
}

// ItemServiceHandler is an implementation of the datapimp.v1.ItemService service.
type ItemServiceHandler interface {
	PushItem(context.Context, *connect_go.Request[v1.PushItemRequest]) (*connect_go.Response[v1.PushItemResponse], error)
	GetItem(context.Context, *connect_go.Request[v1.GetItemRequest]) (*connect_go.Response[v1.GetItemResponse], error)
}

// NewItemServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewItemServiceHandler(svc ItemServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/datapimp.v1.ItemService/PushItem", connect_go.NewUnaryHandler(
		"/datapimp.v1.ItemService/PushItem",
		svc.PushItem,
		opts...,
	))
	mux.Handle("/datapimp.v1.ItemService/GetItem", connect_go.NewUnaryHandler(
		"/datapimp.v1.ItemService/GetItem",
		svc.GetItem,
		opts...,
	))
	return "/datapimp.v1.ItemService/", mux
}

// UnimplementedItemServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedItemServiceHandler struct{}

func (UnimplementedItemServiceHandler) PushItem(context.Context, *connect_go.Request[v1.PushItemRequest]) (*connect_go.Response[v1.PushItemResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("datapimp.v1.ItemService.PushItem is not implemented"))
}

func (UnimplementedItemServiceHandler) GetItem(context.Context, *connect_go.Request[v1.GetItemRequest]) (*connect_go.Response[v1.GetItemResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("datapimp.v1.ItemService.GetItem is not implemented"))
}
