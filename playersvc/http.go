package playersvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

var (
	ErrMissingIds = errors.New("missing ids")
	ErrParsingIds = errors.New("error parsing ids, should be ints")
)

func MakeHTTPHandler(e EndpointSet, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/players").Handler(httptransport.NewServer(
		e.GetAllEndpoint,
		decodeGetAllRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/players").Handler(httptransport.NewServer(
		e.UpdateEndpoint,
		decodeUpdateRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/players").Handler(httptransport.NewServer(
		e.AddEndpoint,
		decodeAddRequest,
		encodeResponse,
		options...,
	))
	return r
}

/** server decode/encode **/
func decodeGetAllRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return getAllRequest{}, nil
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrPlayerDoesNotExist:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

/** client encode/decode **/
func decodeGetAllResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	return decodeMultiResponse(ctx, resp)
}

func decodeAddResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	return decodeSingleResponse(ctx, resp)
}

func decodeUpdateResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	return decodeSingleResponse(ctx, resp)
}

func decodeSingleResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response singleResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func decodeMultiResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response multiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func encodeGetAllRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/players"
	return nil
}

func encodeAddRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/players"
	return encodeRequest(ctx, req, request)
}

func encodeUpdateRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/players"
	return encodeRequest(ctx, req, request)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = io.NopCloser(&buf)
	return nil
}
