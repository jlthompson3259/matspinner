package ticketsvc

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Service interface {
	Get(ctx context.Context, gemIds ...int) ([]Tickets, error)
	Set(ctx context.Context, tickets ...Tickets) ([]Tickets, error)
	Increment(ctx context.Context, gemIds ...int) ([]Tickets, error)
}

type Tickets struct {
	GemId   int `json:"gemId"`
	Tickets int `json:"tickets"`
}

type ticketService struct {
	tickets map[int]int
	logger  log.Logger
}

func NewService(logger log.Logger) Service { return &ticketService{} }

// Get implements Service
func (svc *ticketService) Get(ctx context.Context, gemIds ...int) ([]Tickets, error) {
	tickets := make([]Tickets, len(gemIds))
	for idx, gemId := range gemIds {
		tickets[idx] = Tickets{gemId, svc.tickets[gemId]}
		level.Debug(svc.logger).Log()
	}
	return tickets, nil
}

// Increment implements Service
func (svc *ticketService) Increment(ctx context.Context, gemIds ...int) ([]Tickets, error) {
	tickets := make([]Tickets, len(gemIds))
	for idx, gemId := range gemIds {
		svc.tickets[gemId] += 1
		tickets[idx] = Tickets{gemId, svc.tickets[gemId]}
	}
	return tickets, nil
}

// Set implements Service
func (svc *ticketService) Set(ctx context.Context, tickets ...Tickets) ([]Tickets, error) {
	for _, ticket := range tickets {
		svc.tickets[ticket.GemId] = ticket.Tickets
	}
	return tickets, nil
}
