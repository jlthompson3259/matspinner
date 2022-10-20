package spinsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func MakeHTTPHandler(e EndpointSet, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/spin").Handler(httptransport.NewServer(
		e.SpinEndpoint,
		decodeSpinRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/get-last-spin").Handler(httptransport.NewServer(
		e.GetLastEndpoint,
		decodeGetLastRequest,
		encodeResponse,
		options...,
	))
	return r
}

/** server encode/decode **/
func decodeSpinRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var participantIds []int
	if err := json.NewDecoder(r.Body).Decode(&participantIds); err != nil {
		return nil, err
	}
	return spinRequest{
		ParticipantIds: participantIds,
	}, nil
}

func decodeGetLastRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return getLastRequest{}, nil
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
	case ErrNoSpin, ErrNoTickets:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

/** client encode/decode **/
func decodeResponse(ctx context.Context, resp *http.Response) (interface{}, error) {
	var response response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func encodeSpinRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/spin"
	return encodeRequest(ctx, req, request)
}

func encodeGetLastRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/get-last-spin"
	return nil
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

func encodeIdsQueryString(ids []int) (idStr string) {
	strs := make([]string, len(ids))
	for i, v := range ids {
		strs[i] = fmt.Sprint(v)
	}
	return strings.Join(strs, ",")
}
