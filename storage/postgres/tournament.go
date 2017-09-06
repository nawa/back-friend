package postgres

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/lib/pq"
	"github.com/nawa/back-friend/model"
	"github.com/nawa/back-friend/storage"
)

type tournamentStorage struct {
	db *sql.DB
}

func newTournamentStorage(db *sql.DB) *tournamentStorage {
	tournamentStorage := &tournamentStorage{}
	tournamentStorage.db = db
	return tournamentStorage
}

func (s *tournamentStorage) Create(tournamentID string, deposit int) error {
	stmt, err := s.db.Prepare(`INSERT INTO tournament (id, deposit) VALUES ($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tournamentID, deposit)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			return &storage.IllegalArgumentError{Message: "a tournament with provided id exists"}
		}
		return err
	}
	return nil
}

func (s *tournamentStorage) Join(tournamentID, playerID string, backersID ...string) error {
	return InTransaction(s.db, func(tx *sql.Tx) error {
		deposit, err := joinCheckTournamentStatus(tx, tournamentID)
		if err != nil {
			return err
		}

		err = joinAddParticipants(tx, tournamentID, playerID, backersID...)
		if err != nil {
			return err
		}

		return withdrawPointsForTournament(
			func(query string) (*sql.Stmt, error) {
				return tx.Prepare(query)
			},
			deposit, playerID, backersID...)
	})
}

func joinCheckTournamentStatus(tx *sql.Tx, tournamentID string) (deposit int, err error) {
	rows, err := tx.Query(`SELECT t.deposit, COUNT(w.winner_id) > 0 AS closed
		FROM tournament t LEFT JOIN tournament_winner w ON t.id = w.tournament_id
		WHERE t.id = $1
		GROUP BY t.id`, tournamentID)
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		return 0, &storage.IllegalArgumentError{Message: "tournament not found"}
	}

	var isClosed bool
	err = rows.Scan(&deposit, &isClosed)
	rows.Close()
	if err != nil {
		return 0, err
	}
	if isClosed {
		return 0, &storage.IllegalArgumentError{Message: "tournament with provided id has been closed"}
	}
	return deposit, nil
}

func joinAddParticipants(tx *sql.Tx, tournamentID, playerID string, backersID ...string) error {
	stmt, err := tx.Prepare(
		`WITH conditions AS (SELECT
				  COUNT(w.winner_id) > 0 AS tournament_closed
				FROM tournament t LEFT JOIN tournament_winner w ON t.id = w.tournament_id
				WHERE t.id = $1
				GROUP BY t.id)
		INSERT INTO tournament_participant (tournament_id, player_id, backer_id)
		  SELECT $1, $2, $3 FROM conditions
		  WHERE conditions.tournament_closed = FALSE`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	stmtArgs := [][2]interface{}{}

	if len(backersID) == 0 {
		stmtArgs = append(stmtArgs, [2]interface{}{playerID, nil})
	} else {
		for _, backerID := range backersID {
			stmtArgs = append(stmtArgs, [2]interface{}{playerID, backerID})
		}
	}
	for _, cur := range stmtArgs {
		res, err := stmt.Exec(tournamentID, cur[0], cur[1])
		if err != nil {
			return err
		}
		inserted, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if inserted == 0 {
			return &storage.IllegalArgumentError{Message: "tournament with provided id doesn't exist or has been closed"}
		}
	}
	return nil
}

func (s *tournamentStorage) Result(tournament model.Tournament) error {
	return InTransaction(s.db, func(tx *sql.Tx) error {
		rows, err := tx.Query(`SELECT * FROM tournament WHERE id=$1 LIMIT 1`, tournament.ID)
		if err != nil {
			return err
		}
		if !rows.Next() {
			return &storage.IllegalArgumentError{Message: "tournament not found"}
		}
		rows.Close()

		refundStmt, err := tx.Prepare(`UPDATE player p SET balance = p.balance+$2 WHERE p.id=$1`)
		if err != nil {
			return err
		}
		defer refundStmt.Close()

		for _, winner := range tournament.Winners {
			err = storeTournamentWinner(tx, tournament, winner)
			if err != nil {
				return err
			}
			err = refundToTournamentParticipant(tx, refundStmt, tournament, winner)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func storeTournamentWinner(tx *sql.Tx, tournament model.Tournament, winner model.Winner) error {
	stmt, err := tx.Prepare(`INSERT INTO tournament_winner (tournament_id, winner_id, prize) VALUES ($1, $2, $3)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tournament.ID, winner.ID, winner.Prize)
	return err
}

func refundToTournamentParticipant(tx *sql.Tx, refundStmt *sql.Stmt, tournament model.Tournament, winner model.Winner) error {
	rows, err := tx.Query(`SELECT backer_id FROM tournament_participant WHERE tournament_id=$1 AND player_id=$2`,
		tournament.ID, winner.ID)
	if err != nil {
		return err
	}

	var backersID = []string{}
	winnerExists := false
	for rows.Next() {
		var backerID string
		err = rows.Scan(&backerID)
		if err != nil {
			rows.Close()
			return err
		}
		if backerID != "" {
			backersID = append(backersID, backerID)
		}
		winnerExists = true
	}
	rows.Close()
	if !winnerExists {
		return &storage.IllegalArgumentError{
			Message: fmt.Sprintf("winner with id '%v' not found", winner.ID),
		}
	}

	part := int(math.Ceil(float64(winner.Prize) / float64(len(backersID)+1)))
	_, err = refundStmt.Exec(winner.ID, part)
	if err != nil {
		return err
	}

	if len(backersID) > 0 {
		for _, backerID := range backersID {
			_, err = refundStmt.Exec(backerID, part)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
