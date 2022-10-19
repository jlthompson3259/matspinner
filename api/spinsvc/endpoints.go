package spinsvc

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type EndpointSet struct {
	SpinEndpoint    endpoint.Endpoint
	GetLastEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(svc Service) EndpointSet {
	return EndpointSet{
		SpinEndpoint:    MakeSpinEndpoint(svc),
		GetLastEndpoint: MakeGetLastEndpoint(svc),
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
		SpinEndpoint:    httptransport.NewClient("POST", tgt, encodeSpinRequest, decodeResponse, options...).Endpoint(),
		GetLastEndpoint: httptransport.NewClient("GET", tgt, encodeGetLastRequest, decodeResponse, options...).Endpoint(),
	}, nil
}

func (e *EndpointSet) Spin(ctx context.Context, participantids []int) (SpinResult, error) {
	request := spinRequest{ParticipantIds: participantids}
	r, err := e.SpinEndpoint(ctx, request)
	if err != nil {
		return SpinResult{}, err
	}
	resp := r.(response)
	return resp.Result, nil
}

func (e *EndpointSet) GetLast(ctx context.Context) (SpinResult, error) {
	request := getLastRequest{}
	r, err := e.GetLastEndpoint(ctx, request)
	if err != nil {
		return SpinResult{}, err
	}
	resp := r.(response)
	return resp.Result, nil
}

func MakeSpinEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(spinRequest)
		tickets, err := svc.Spin(ctx, req.ParticipantIds)
		return response{tickets, err}, nil
	}
}

func MakeGetLastEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tickets, err := svc.GetLast(ctx)
		return response{tickets, err}, nil
	}
}

type spinRequest struct {
	ParticipantIds []int
}

type getLastRequest struct {
}

type response struct {
	Result SpinResult `json:"result,omitempty"`
	Error  error      `json:"error,omitempty"`
}

func (r response) error() error {
	return r.Error
}
