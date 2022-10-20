package ticketsvc

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type EndpointSet struct {
	GetEndpoint       endpoint.Endpoint
	SetEndpoint       endpoint.Endpoint
	IncrementEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(svc Service) EndpointSet {
	return EndpointSet{
		GetEndpoint:       MakeGetEndpoint(svc),
		SetEndpoint:       MakeSetEndpoint(svc),
		IncrementEndpoint: MakeIncrementEndpoint(svc),
	}
}

func MakeClientEndpoints(instance string) (EndpointSet, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return EndpointSet{}, err
	}

	tgt.Path = ""

	options := []httptransport.ClientOption{}

	return EndpointSet{
		GetEndpoint:       httptransport.NewClient("GET", tgt, encodeGetRequest, decodeResponse, options...).Endpoint(),
		SetEndpoint:       httptransport.NewClient("PUT", tgt, encodeSetRequest, decodeResponse, options...).Endpoint(),
		IncrementEndpoint: httptransport.NewClient("POST", tgt, encodeIncrementRequest, decodeResponse, options...).Endpoint(),
	}, nil
}

func (e *EndpointSet) Get(ctx context.Context, ids ...int) ([]Tickets, error) {
	request := getRequest{Ids: ids}
	r, err := e.GetEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := r.(response)
	return resp.Tickets, nil
}

func (e *EndpointSet) Set(ctx context.Context, tickets ...Tickets) ([]Tickets, error) {
	request := setRequest{Tickets: tickets}
	r, err := e.SetEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := r.(response)
	return resp.Tickets, nil
}

func (e *EndpointSet) Increment(ctx context.Context, ids ...int) ([]Tickets, error) {
	request := incrementRequest{Ids: ids}
	r, err := e.IncrementEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := r.(response)
	return resp.Tickets, nil
}

func MakeGetEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		tickets, err := svc.Get(ctx, req.Ids...)
		return response{tickets, err}, nil
	}
}

func MakeSetEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(setRequest)
		tickets, err := svc.Set(ctx, req.Tickets...)
		return response{tickets, err}, nil
	}
}

func MakeIncrementEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(incrementRequest)
		tickets, err := svc.Increment(ctx, req.Ids...)
		return response{tickets, err}, nil
	}
}

type getRequest struct {
	Ids []int
}

type incrementRequest struct {
	Ids []int
}

type setRequest struct {
	Tickets []Tickets
}

type response struct {
	Tickets []Tickets `json:"tickets,omitempty"`
	Err     error     `json:"err,omitempty"`
}

func (r response) error() error { return r.Err }
