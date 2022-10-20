package ticketsvc

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Service interface {
	Get(ctx context.Context, ids ...int) ([]Tickets, error)
	Set(ctx context.Context, tickets ...Tickets) ([]Tickets, error)
	Increment(ctx context.Context, ids ...int) ([]Tickets, error)
}

type Tickets struct {
	Id      int `json:"id"`
	Tickets int `json:"tickets"`
}

func (t Tickets) String() string {
	return fmt.Sprintf("{id: %v, tickets: %v}", t.Id, t.Tickets)
}

type ticketService struct {
	tickets map[int]int
	logger  log.Logger
}

func NewService(logger log.Logger) Service {
	return &ticketService{
		logger:  logger,
		tickets: make(map[int]int),
	}
}

// Get implements Service
func (svc *ticketService) Get(ctx context.Context, ids ...int) ([]Tickets, error) {
	tickets := make([]Tickets, len(ids))
	for idx, id := range ids {
		t := svc.tickets[id]
		tickets[idx] = Tickets{id, t}
		level.Debug(svc.logger).Log("debug", "get tickets", "id", id, "tickets", t)
	}
	return tickets, nil
}

// Increment implements Service
func (svc *ticketService) Increment(ctx context.Context, ids ...int) ([]Tickets, error) {
	tickets := make([]Tickets, len(ids))
	for idx, id := range ids {
		svc.tickets[id] += 1
		t := svc.tickets[id]
		tickets[idx] = Tickets{id, t}
		level.Debug(svc.logger).Log("debug", "increment tickets", "id", id, "tickets", t)
	}
	return tickets, nil
}

// Set implements Service
func (svc *ticketService) Set(ctx context.Context, tickets ...Tickets) ([]Tickets, error) {
	for _, t := range tickets {
		ticks := t.Tickets
		if ticks <= 0 {
			delete(svc.tickets, t.Id)
		} else {
			svc.tickets[t.Id] = t.Tickets
		}
		level.Debug(svc.logger).Log("debug", "set tickets", "id", t.Id, "tickets", t.Tickets)
	}
	return tickets, nil
}
