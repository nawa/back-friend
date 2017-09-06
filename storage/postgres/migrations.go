package postgres

import "github.com/rubenv/sql-migrate"

var migrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		{
			Id: "1_init_schema",
			Up: []string{`
CREATE TABLE player
(
  id      VARCHAR PRIMARY KEY NOT NULL,
  balance INT DEFAULT 0       NOT NULL
);

CREATE TABLE tournament
(
  id      VARCHAR PRIMARY KEY NOT NULL,
  deposit INT                 NOT NULL
);

CREATE TABLE tournament_participant
(
  tournament_id VARCHAR NOT NULL,
  player_id     VARCHAR NOT NULL,
  backer_id     VARCHAR,
  CONSTRAINT tournament_participant_tournament_id_player_id_backer_id_pk UNIQUE (tournament_id, player_id, backer_id),
  CONSTRAINT tournament_participant_tournament_id_fk FOREIGN KEY (tournament_id) REFERENCES tournament (id),
  CONSTRAINT tournament_participant_player_id_fk FOREIGN KEY (player_id) REFERENCES player (id),
  CONSTRAINT tournament_participant_backer_id_fk FOREIGN KEY (backer_id) REFERENCES player (id)
);

CREATE TABLE tournament_winner
(
  tournament_id VARCHAR NOT NULL,
  winner_id     VARCHAR NOT NULL,
  prize         INT     NOT NULL,
  CONSTRAINT tournament_winner_tournament_id_fk FOREIGN KEY (tournament_id) REFERENCES tournament (id),
  CONSTRAINT tournament_winner_player_id_fk FOREIGN KEY (winner_id) REFERENCES player (id)
);`,
			},
			Down: []string{`
DROP TABLE tournament_winner;
DROP TABLE tournament_participant;
DROP TABLE tournament;
DROP TABLE player;`,
			},
		},
	},
}
