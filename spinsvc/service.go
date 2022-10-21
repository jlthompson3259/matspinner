package spinsvc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-kit/log"

	"github.com/jlthompson3259/matspinner/ticketsvc"
)

var (
	ErrNoParticipants = errors.New("no participants")
	ErrNoTickets      = errors.New("none of the participants have tickets")
	ErrNoSpin         = errors.New("no spin yet to return")
)

type Service interface {
	Spin(ctx context.Context, participantIds []int) (SpinResult, error)
	SpinUnweighted(ctx context.Context, particantIds []int) (SpinResult, error)
	GetLast(ctx context.Context) (SpinResult, error)
}

type SpinResult struct {
	ParticipantIds []int `json:"participantIds"`
	WinnerId       int   `json:"winnerId"`
}

func (t SpinResult) String() string {
	return fmt.Sprintf("{participants: %v, winner: %v}", t.ParticipantIds, t.WinnerId)
}

type spinService struct {
	logger        log.Logger
	ticketService ticketsvc.Service
	lastSpin      *SpinResult
}

func NewService(logger log.Logger, ticketService ticketsvc.Service) Service {
	return &spinService{
		logger:        logger,
		ticketService: ticketService,
	}
}

func (s *spinService) Spin(ctx context.Context, participantIds []int) (SpinResult, error) {
	return s.spinUsingTicketFunction(ctx, participantIds, func(tickets []ticketsvc.Tickets) []ticketsvc.Tickets {
		return tickets
	})
}

func (s *spinService) SpinUnweighted(ctx context.Context, participantIds []int) (SpinResult, error) {
	return s.spinUsingTicketFunction(ctx, participantIds, func(tickets []ticketsvc.Tickets) []ticketsvc.Tickets {
		// give each participant only 1 ticket for the spin
		ret := make([]ticketsvc.Tickets, len(tickets))
		for i, v := range tickets {
			ret[i] = ticketsvc.Tickets{Id: v.Id, Tickets: 1}
		}
		return ret
	})
}

func (s *spinService) GetLast(ctx context.Context) (SpinResult, error) {
	if s.lastSpin == nil {
		return SpinResult{}, ErrNoSpin
	}

	return *s.lastSpin, nil
}

func (s *spinService) spinUsingTicketFunction(
	ctx context.Context,
	participantIds []int,
	ticketFunc func(tickets []ticketsvc.Tickets) []ticketsvc.Tickets,
) (SpinResult, error) {

	if len(participantIds) <= 0 {
		return SpinResult{}, ErrNoParticipants
	}

	tickets, err := s.ticketService.Increment(ctx, participantIds...)
	if err != nil {
		return SpinResult{}, err
	}

	tickets = ticketFunc(tickets)

	winner, err := chooseRandomWinner(tickets)
	if err != nil {
		return SpinResult{}, err
	}

	result := SpinResult{
		ParticipantIds: participantIds,
		WinnerId:       winner,
	}
	s.lastSpin = &result

	_, err = s.ticketService.Set(ctx, ticketsvc.Tickets{Id: winner, Tickets: 0})
	if err != nil {
		return result, err
	}

	return result, nil
}

func chooseRandomWinner(tickets []ticketsvc.Tickets) (winnerId int, err error) {
	ticketSum := 0
	for _, v := range tickets {
		ticketSum += v.Tickets
	}

	if ticketSum <= 0 {
		return 0, ErrNoTickets
	}

	rand.Seed(time.Now().UnixNano())
	winningTicket := rand.Intn(ticketSum)

	for _, v := range tickets {
		winningTicket -= v.Tickets
		if winningTicket < 0 {
			winnerId = v.Id
			break
		}
	}

	return winnerId, nil
}
