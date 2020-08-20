package db

import (
	"github.com/mellena1/RSC-Spreadsheet-API/data/models"
	log "github.com/sirupsen/logrus"
)

func (db *DB) GetAllTeams() ([]models.Team, error) {
	rows, err := db.sqlDB.Query(`
		SELECT team_id, name, franchise, conference, tier, division FROM team;
	`)
	if err != nil {
		log.Errorf("Error getting all teams from db: %s", err)
		return nil, err
	}
	defer rows.Close()

	teams := []models.Team{}
	for rows.Next() {
		team := models.Team{}
		err := rows.Scan(&team.TeamID, &team.Name, &team.Franchise, &team.Conference, &team.Tier, &team.Division)
		if err != nil {
			log.Errorf("Error scanning a team: %s", err)
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}
