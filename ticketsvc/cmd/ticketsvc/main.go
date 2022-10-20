package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/jlthompson3259/matspinner/ticketsvc"
)

const (
	defaultHttpPort = "8085"
)

func main() {
	var (
		httpAddr = net.JoinHostPort("0.0.0.0", envString("HTTP_PORT", defaultHttpPort))
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = level.NewFilter(logger, level.AllowInfo(), level.SquelchNoLevel(true))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	var service ticketsvc.Service
	{
		service = ticketsvc.NewService(log.With(logger, "component", "service"))
		service = ticketsvc.LoggingMiddleware(log.With(logger, "component", "loggingMiddleware"))(service)
	}

	var (
		endpoints   = ticketsvc.MakeServerEndpoints(service)
		httpHandler = ticketsvc.MakeHTTPHandler(endpoints, log.With(logger, "component", "http"))
	)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", httpAddr)
		errs <- http.ListenAndServe(httpAddr, httpHandler)
	}()

	level.Error(logger).Log("exit", <-errs)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
