package spinsvc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/jlthompson3259/matspinner/ticketsvc"
)

var (
	ErrNoTickets = errors.New("no participants or none of the participants have tickets")
	ErrNoSpin    = errors.New("no spin yet to return")
)

type Service interface {
	Spin(ctx context.Context, participantIds []int) (SpinResult, error)
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
	tickets, err := s.ticketService.Increment(ctx, participantIds...)
	if err != nil {
		return SpinResult{}, err
	}

	level.Info(s.logger).Log("tickets", fmt.Sprintf("%v", tickets))
	ticketSum := 0
	for _, v := range tickets {
		ticketSum += v.Tickets
	}

	if ticketSum <= 0 {
		return SpinResult{}, ErrNoTickets
	}

	rand.Seed(time.Now().UnixNano())
	winningTicket := rand.Intn(ticketSum)

	var winner int
	for _, v := range tickets {
		winningTicket -= v.Tickets
		if winningTicket < 0 {
			winner = v.Id
			break
		}
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

func (s *spinService) GetLast(ctx context.Context) (SpinResult, error) {
	if s.lastSpin == nil {
		return SpinResult{}, ErrNoSpin
	}

	return *s.lastSpin, nil
}
