package main

import (
	"net"
	"os"
	"ticketsvc"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	defaultHttpPort = "1337"
)

func main() {
	var (
		httpAddr = net.JoinHostPort("localhost", envString("HTTP_PORT", defaultHttpPort))
	)

	var logger log.Logger
	{
		logger = level.NewFilter()
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	var service ticketsvc.Service
	{
		service = ticketsvc.NewService(log.With(logger, "component", "service"))
		service = ticketsvc.LoggingMiddleware()
	}

	var (
		endpoints   = ticketsvc.MakeServerEndpoints(service)
		httpHandler = ticketsvc.MakeHTTPHandler(service, log.With(logger, "component", "http"))
	)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
