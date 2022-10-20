package playersvc

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type EndpointSet struct {
	GetAllEndpoint endpoint.Endpoint
	AddEndpoint    endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(svc Service) EndpointSet {
	return EndpointSet{
		GetAllEndpoint: MakeGetAllEndpoint(svc),
		AddEndpoint:    MakeAddEndpoint(svc),
		UpdateEndpoint: MakeUpdateEndpoint(svc),
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
		GetAllEndpoint: httptransport.NewClient("GET", tgt, encodeGetAllRequest, decodeGetAllResponse, options...).Endpoint(),
		AddEndpoint:    httptransport.NewClient("POST", tgt, encodeAddRequest, decodeAddResponse, options...).Endpoint(),
		UpdateEndpoint: httptransport.NewClient("PUT", tgt, encodeUpdateRequest, decodeUpdateResponse, options...).Endpoint(),
	}, nil
}

func (e *EndpointSet) GetAll(ctx context.Context) ([]Player, error) {
	request := getAllRequest{}
	r, err := e.GetAllEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := r.(multiResponse)
	return resp.Players, nil
}

func (e *EndpointSet) Add(ctx context.Context, name string) (Player, error) {
	request := addRequest{Name: name}
	r, err := e.AddEndpoint(ctx, request)
	if err != nil {
		return Player{}, err
	}
	resp := r.(singleResponse)
	return resp.Player, nil
}

func (e *EndpointSet) Increment(ctx context.Context, player Player) (Player, error) {
	request := updateRequest{Player: player}
	r, err := e.UpdateEndpoint(ctx, request)
	if err != nil {
		return Player{}, err
	}
	resp := r.(singleResponse)
	return resp.Player, nil
}

func MakeGetAllEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		players, err := svc.GetAll(ctx)
		return multiResponse{players, err}, nil
	}
}

func MakeAddEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		player, err := svc.Add(ctx, req.Name)
		return singleResponse{player, err}, nil
	}
}

func MakeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRequest)
		player, err := svc.Update(ctx, req.Player)
		return singleResponse{player, err}, nil
	}
}

type getAllRequest struct {
}

type updateRequest struct {
	Player Player `json:"player,omitempty"`
}

type addRequest struct {
	Name string `json:"name,omitempty"`
}

type singleResponse struct {
	Player Player `json:"player,omitempty"`
	Err    error  `json:"err,omitempty"`
}

func (r singleResponse) error() error { return r.Err }

type multiResponse struct {
	Players []Player `json:"players,omitempty"`
	Err     error    `json:"err,omitempty"`
}

func (r multiResponse) error() error { return r.Err }
