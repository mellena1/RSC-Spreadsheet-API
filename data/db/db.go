package db

import (
	"database/sql"

	"github.com/mellena1/RSC-Spreadsheet-API/data/models"
	"github.com/mellena1/RSC-Spreadsheet-API/data/sheets"
	log "github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

type Datastore interface {
	GetAllTeams(GetAllTeamsQuery) ([]models.Team, error)
}

type DB struct {
	sqlDB *sql.DB

	teamStandingsUpdater sheets.TeamStandingsRetriever
}

func NewDB(connStr string, teamStandingsSheet sheets.TeamStandingsRetriever) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	newdb := &DB{
		sqlDB: db,

		teamStandingsUpdater: teamStandingsSheet,
	}

	if err = newdb.makeTablesIfNotExist(); err != nil {
		return nil, err
	}

	if err = newdb.fillTeamData(); err != nil {
		return nil, err
	}

	return newdb, nil
}

func (db *DB) Close() error {
	return db.sqlDB.Close()
}

func (db *DB) makeTablesIfNotExist() error {
	tx, err := db.sqlDB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS team (
			team_id SERIAL PRIMARY KEY,
			name text NOT NULL,
			franchise text NOT NULL,
			conference text NOT NULL,
			tier text NOT NULL,
			division text,
			UNIQUE(name, franchise, tier)
		);
	`)
	if err != nil {
		log.Errorf("Failed to make team table: %v", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) fillTeamData() error {
	teamData, err := db.teamStandingsUpdater.GetTeamStandingsFromSheet()
	if err != nil {
		return err
	}

	tx, err := db.sqlDB.Begin()
	if err != nil {
		return err
	}

	for _, t := range teamData {
		_, err = tx.Exec(`
			INSERT INTO team (name, franchise, conference, tier, division) 
			VALUES($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING;
		`, t.Team.Name, t.Team.Franchise, t.Team.Conference, t.Team.Tier, t.Team.Division)
		if err != nil {
			log.Errorf("Failed to insert team into team table: %v", err)
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
