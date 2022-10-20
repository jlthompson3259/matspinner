package playersvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/log"
)

var (
	ErrPlayerDoesNotExist = errors.New("player does not exist")
)

type Service interface {
	Add(ctx context.Context, name string) (Player, error)
	GetAll(ctx context.Context) ([]Player, error)
	Update(ctx context.Context, player Player) (Player, error)
}

type Player struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (t Player) String() string {
	return fmt.Sprintf("{id: %v, name: %v}", t.Id, t.Name)
}

type playerService struct {
	players map[int]string
	logger  log.Logger
}

func NewService(logger log.Logger) Service {
	return &playerService{
		players: make(map[int]string),
		logger:  logger,
	}
}

func (s *playerService) Add(ctx context.Context, name string) (Player, error) {
	idx := len(s.players)
	s.players[idx] = name
	return Player{Id: idx, Name: name}, nil
}

func (s *playerService) GetAll(ctx context.Context) (players []Player, err error) {
	players = make([]Player, len(s.players))
	for id, name := range s.players {
		players = append(players, Player{Id: id, Name: name})
	}
	return
}

func (s *playerService) Update(ctx context.Context, player Player) (Player, error) {

	if _, ok := s.players[player.Id]; !ok {
		return Player{}, ErrPlayerDoesNotExist
	}
	s.players[player.Id] = player.Name
	return player, nil
}
