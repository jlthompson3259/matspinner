package ticketsvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
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

func MakeGetEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)
		tickets, err := svc.Get(ctx, req.GemIds...)
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
		tickets, err := svc.Increment(ctx, req.GemIds...)
		return response{tickets, err}, nil
	}
}

type getRequest struct {
	GemIds []int
}

type incrementRequest struct {
	GemIds []int
}

type setRequest struct {
	Tickets []Tickets
}

type response struct {
	Tickets []Tickets `json:"tickets,omitempty"`
	Error   error     `json:"error,omitempty"`
}

func (r response) error() error {
	return r.Error
}
