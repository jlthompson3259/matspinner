package ticketsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

	r.Methods("GET").Path("/tickets").Handler(httptransport.NewServer(
		e.GetEndpoint,
		decodeGetRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/tickets").Handler(httptransport.NewServer(
		e.SetEndpoint,
		decodeSetRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/tickets/increment").Handler(httptransport.NewServer(
		e.IncrementEndpoint,
		decodeIncrementRequest,
		encodeResponse,
		options...,
	))
	return r
}

/** server decode/encode **/
func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	if !q.Has("ids") {
		return nil, ErrMissingIds
	}
	ids, err := decodeIdsQueryString(q.Get("ids"))
	if err != nil {
		return nil, ErrParsingIds
	}
	return getRequest{Ids: ids}, nil
}

func decodeSetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request setRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeIncrementRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request incrementRequest
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

func decodeIdsQueryString(idStr string) (ids []int, err error) {
	strs := strings.Split(idStr, ",")
	ids = make([]int, len(strs))
	for i, v := range strs {
		ids[i], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}
	return ids, nil
}

func codeFrom(err error) int {
	switch err {
	case ErrMissingIds, ErrParsingIds:
		return http.StatusBadRequest
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

func encodeGetRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(getRequest)
	var qStr string
	{
		qStr = encodeIdsQueryString(r.Ids)
		qStr = url.QueryEscape(qStr)
	}
	req.URL.Path = "/tickets"
	req.URL.Query().Add("ids", qStr)
	return nil
}

func encodeSetRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/tickets"
	return encodeRequest(ctx, req, request)
}

func encodeIncrementRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.URL.Path = "/tickets/increment"
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

func encodeIdsQueryString(ids []int) (idStr string) {
	strs := make([]string, len(ids))
	for i, v := range ids {
		strs[i] = fmt.Sprint(v)
	}
	return strings.Join(strs, ",")
}
