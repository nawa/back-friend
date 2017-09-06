package storage

import (
	"errors"

	"github.com/nawa/back-friend/model"
)

type Storage interface {
	PlayerStorage() PlayerStorage
	TournamentStorage() TournamentStorage
	Reset() error
}

type PlayerStorage interface {
	Get(playerID string) (*model.Player, error)
	Take(playerID string, points int) error
	FundOrCreate(playerID string, points int) error
}

type TournamentStorage interface {
	Create(tournamentID string, deposit int) error
	Join(tournamentID, playerID string, backersID ...string) error
	Result(tournament model.Tournament) error
}

var ErrNotFound = errors.New("entity not found")

type IllegalArgumentError struct {
	Message string
}

func (e *IllegalArgumentError) Error() string {
	return e.Message
}
