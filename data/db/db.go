package db

import (
	"database/sql"

	"github.com/mellena1/RSC-Spreadsheet-API/data/sheets"
	"github.com/mellena1/RSC-Spreadsheet-API/models"
)

type Datastore interface {
	AllTeamStandings() ([]models.TeamStanding, error)
}

type DB struct {
	sqlDB *sql.DB

	teamStandingsUpdater sheets.TeamStandingsRetriever
}

func NewDB(dataSourceName string, teamStandingsSheet sheets.TeamStandingsRetriever) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
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

	return newdb, nil
}

func (d *DB) makeTablesIfNotExist() error {
	tx, err := d.sqlDB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS teamstanding (
			teamID text NOTNULL

		);
	`)
	if err != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}
