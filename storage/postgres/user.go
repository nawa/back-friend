package postgres

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/nawa/back-friend/model"
	"github.com/nawa/back-friend/storage"
)

type playerStorage struct {
	db *sql.DB
}

func newPlayerStorage(db *sql.DB) *playerStorage {
	playerStorage := &playerStorage{}
	playerStorage.db = db
	return playerStorage
}

func (s *playerStorage) Get(playerID string) (*model.Player, error) {
	var player = new(model.Player)
	rows, err := s.db.Query(`SELECT id, balance FROM player WHERE id = $1`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, storage.ErrNotFound
	}
	err = rows.Scan(&player.ID, &player.Balance)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (s *playerStorage) Take(playerID string, points int) error {
	if points <= 0 {
		return &storage.IllegalArgumentError{Message: "Incorrect points value to take"}
	}
	q := `UPDATE player p
		SET balance = CASE
		              WHEN p.balance >= $2
		                THEN p.balance - $2
		              ELSE p.balance
		              END
		FROM player old_p
		WHERE p.id = old_p.id AND p.id = $1
		RETURNING old_p.balance`

	stmt, err := s.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(playerID, points)
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return storage.ErrNotFound
	}
	var oldBalance int
	err = rows.Scan(&oldBalance)
	if err != nil {
		return err
	}

	if oldBalance < points {
		return &storage.IllegalArgumentError{Message: "Not enough points on balance to take"}
	}

	return nil
}

func (s *playerStorage) FundOrCreate(playerID string, points int) error {
	stmt, err := s.db.Prepare(
		`INSERT INTO player (id, balance) VALUES ($1, $2)
			ON CONFLICT (id)
			DO UPDATE SET balance = player.balance + EXCLUDED.balance`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(playerID, points)
	return err
}

func withdrawPointsForTournament(stmtCreator func(query string) (*sql.Stmt, error),
	deposit int, playerID string, backersID ...string) error {
	var part int
	if len(backersID) == 0 {
		part = deposit
	} else {
		//rake with rounding :)
		part = int(math.Ceil(float64(deposit) / float64(len(backersID)+1)))
	}

	stmt, err := stmtCreator(`UPDATE player p
		SET balance = CASE
		              WHEN p.balance >= $2
		                THEN p.balance - $2
		              ELSE p.balance
		              END
		FROM player old_p
		WHERE p.id = old_p.id AND p.id = $1
		RETURNING old_p.balance`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmtArgs := []interface{}{}
	stmtArgs = append(stmtArgs, playerID)
	if len(backersID) > 0 {
		for _, backerID := range backersID {
			stmtArgs = append(stmtArgs, backerID)
		}
	}

	for _, cur := range stmtArgs {
		rows, err := stmt.Query(cur, part)
		if err != nil {
			return err
		}

		if !rows.Next() {
			return storage.ErrNotFound
		}
		var oldBalance int
		err = rows.Scan(&oldBalance)
		if err != nil {
			return err
		}

		rows.Close()

		if oldBalance < part {
			return &storage.IllegalArgumentError{
				Message: fmt.Sprintf("Player with id '%v' has not enough points on balance to join", cur),
			}
		}
	}

	return nil
}
