package postgres

import (
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq" //import from postgres package
	"github.com/nawa/back-friend/config"
	"github.com/nawa/back-friend/storage"
	"github.com/rubenv/sql-migrate"
)

type postgresStorage struct {
	config            config.SQLDb
	db                *sql.DB
	playerStorage     storage.PlayerStorage
	tournamentStorage storage.TournamentStorage
}

func NewPostgresStorage(config config.SQLDb) (storage.Storage, error) {
	postgresStorage := &postgresStorage{config: config}
	err := postgresStorage.init()
	if err != nil {
		return nil, err
	}
	postgresStorage.injectStorages()
	return postgresStorage, nil
}

func (s *postgresStorage) PlayerStorage() storage.PlayerStorage {
	return s.playerStorage
}

func (s *postgresStorage) TournamentStorage() storage.TournamentStorage {
	return s.tournamentStorage
}

func (s *postgresStorage) Reset() error {
	_, err := s.db.Exec(`
		DROP TABLE tournament_winner CASCADE;
		DROP TABLE tournament_participant CASCADE;
		DROP TABLE tournament CASCADE;
		DROP TABLE player CASCADE;
		DROP TABLE gorp_migrations CASCADE;`)
	if err != nil {
		return err
	}

	return s.applyAllMigrations()
}

func (s *postgresStorage) init() (err error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		s.config.Host, s.config.Port, s.config.Name, s.config.User, s.config.Password)
	s.db, err = sql.Open(s.config.Type, dataSourceName)
	if err != nil {
		return fmt.Errorf("Fail to connect to database: %v", err)
	}

	return s.applyAllMigrations()
}

func (s *postgresStorage) applyAllMigrations() error {
	n, err := migrate.Exec(s.db, s.config.Type, migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("Fail to perform migrations: %v", err)
	}
	log.Debug("Count of applied migrations - %v", n)
	return nil
}

func (s *postgresStorage) injectStorages() {
	s.playerStorage = newPlayerStorage(s.db)
	s.tournamentStorage = newTournamentStorage(s.db)
}

func InTransaction(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()
	err = fn(tx)
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return fmt.Errorf("failed to rollback transaction: %v; source error: %v", e.Error(), err.Error())
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err.Error())
	}
	return nil
}
