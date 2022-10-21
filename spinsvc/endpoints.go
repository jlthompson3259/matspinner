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
	request := spinRequest{ParticipantIds: participantids, Unweighted: false}
	r, err := e.SpinEndpoint(ctx, request)
	if err != nil {
		return SpinResult{}, err
	}
	resp := r.(response)
	return resp.Result, nil
}

func (e *EndpointSet) SpinUnweighted(ctx context.Context, participantids []int) (SpinResult, error) {
	request := spinRequest{ParticipantIds: participantids, Unweighted: true}
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
	return func(ctx context.Context, request interface{}) (r interface{}, err error) {
		req := request.(spinRequest)
		var result SpinResult
		if req.Unweighted {
			result, err = svc.SpinUnweighted(ctx, req.ParticipantIds)
		} else {
			result, err = svc.Spin(ctx, req.ParticipantIds)
		}
		return response{result, err}, nil
	}
}

func MakeGetLastEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		result, err := svc.GetLast(ctx)
		return response{result, err}, nil
	}
}

type spinRequest struct {
	ParticipantIds []int `json:"participantIds"`
	Unweighted     bool  `json:"unweighted"`
}

type getLastRequest struct {
}

type response struct {
	Result SpinResult `json:"result,omitempty"`
	Err    error      `json:"err,omitempty"`
}

func (r response) error() error { return r.Err }
