package ticketsvc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
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

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
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

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	if !q.Has("ids") {
		return nil, ErrMissingIds
	}
	idStrs := strings.Split(q.Get("ids"), ",")
	ids, err := parseIdStringArray(idStrs)
	if err != nil {
		return nil, ErrParsingIds
	}
	return getRequest{GemIds: ids}, nil
}

func decodeSetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var tickets []Tickets
	if err := json.NewDecoder(r.Body).Decode(&tickets); err != nil {
		return nil, err
	}
	return setRequest{
		Tickets: tickets,
	}, nil
}

func decodeIncrementRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var ids []int
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		return nil, err
	}
	return incrementRequest{
		GemIds: ids,
	}, nil
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

func parseIdStringArray(idStrs []string) (ids []int, err error) {
	ids = make([]int, len(idStrs))
	for i, v := range idStrs {
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
