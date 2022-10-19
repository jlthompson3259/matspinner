package ticketsvc

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

type ServiceMiddleware func(Service) Service

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(service Service) Service {
		return &loggingMiddleware{
			next:   service,
			logger: logger,
		}
	}
}

func (mw *loggingMiddleware) Get(ctx context.Context, ids ...int) (t []Tickets, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Get", "ids", ids, "duration", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Get(ctx, ids...)
}

func (mw *loggingMiddleware) Set(ctx context.Context, tickets ...Tickets) (t []Tickets, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Set", "data", tickets, "duration", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Set(ctx, tickets...)
}

func (mw *loggingMiddleware) Increment(ctx context.Context, ids ...int) (t []Tickets, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Increment", "ids", ids, "duration", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Increment(ctx, ids...)
}
