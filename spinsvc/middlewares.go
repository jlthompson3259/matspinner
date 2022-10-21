package spinsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
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

func (mw *loggingMiddleware) Spin(ctx context.Context, participantIds []int) (res SpinResult, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "Spin", "result", fmt.Sprintf("%v", res), "duration", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Spin(ctx, participantIds)
}

func (mw *loggingMiddleware) SpinUnweighted(ctx context.Context, participantIds []int) (res SpinResult, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "SpinUnweighted", "result", fmt.Sprintf("%v", res), "duration", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.SpinUnweighted(ctx, participantIds)
}

func (mw *loggingMiddleware) GetLast(ctx context.Context) (res SpinResult, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "GetLast", "result", fmt.Sprintf("%v", res), "duration", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetLast(ctx)
}
