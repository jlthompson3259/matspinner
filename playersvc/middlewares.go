package playersvc

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

func (mw *loggingMiddleware) Add(ctx context.Context, name string) (p Player, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "Add", "player", fmt.Sprintf("%v", p), "duration", time.Since(begin), "err", err)
	}(time.Now())
	p, err = mw.next.Add(ctx, name)
	return
}

func (mw *loggingMiddleware) GetAll(ctx context.Context) (p []Player, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "GetAll", "players", fmt.Sprintf("%v", p), "duration", time.Since(begin), "err", err)
	}(time.Now())
	p, err = mw.next.GetAll(ctx)
	return
}

func (mw *loggingMiddleware) Update(ctx context.Context, player Player) (p Player, err error) {
	defer func(begin time.Time) {
		level.Info(mw.logger).Log("method", "Update", "player", fmt.Sprintf("%v", player), "duration", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Update(ctx, player)
}
