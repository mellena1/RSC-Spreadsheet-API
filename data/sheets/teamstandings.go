package sheets

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mellena1/RSC-Spreadsheet-API/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const TEAMSTANDINGSHEADERS = 2

type TeamStandingsRetriever interface {
	GetTeamStandingsFromSheet() ([]models.TeamStanding, error)
}

type TeamStandingsSheet struct {
	sheetsService *sheets.Service
	spreadsheetID string
	sheetName     string
}

func NewTeamStandingsSheet(ctx context.Context, spreadsheetID, sheetName, apiKey string) (*TeamStandingsSheet, error) {
	svc, err := sheets.NewService(ctx, option.WithAPIKey(apiKey))
	return &TeamStandingsSheet{
		sheetsService: svc,
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
	}, err
}

func (t TeamStandingsSheet) GetTeamStandingsFromSheet() ([]models.TeamStanding, error) {
	result, err := t.sheetsService.Spreadsheets.Values.Get(t.spreadsheetID, t.sheetName).Do()
	if err != nil {
		return nil, err
	}

	rows := result.Values[TEAMSTANDINGSHEADERS:]

	standings := make([]models.TeamStanding, 0, len(rows))
	for i, row := range rows {
		standing, err := rowToTeamStanding(row)
		if err != nil {
			log.Errorf("Row %d failed to be converted: %v", i, row)
			continue
		}
		standings = append(standings, standing)
	}

	return standings, nil
}

func setTeamStandingValBasedOnColumn(t *models.TeamStanding, colIndex int, val interface{}) error {
	valStr, err := assertToString(val)
	if err != nil {
		return err
	}

	switch colIndex {
	case 0:
		t.Team.Tier = valStr
	case 1:
		t.Team.Franchise = valStr
	case 2:
		t.Team.Name = valStr
	case 3:
		t.Team.Conference = valStr
	case 4:
		if valStr == "N/A" {
			t.Team.Division = nil
		} else {
			t.Team.Division = &valStr
		}
	case 7:
		wins, err := strconv.Atoi(valStr)
		if err != nil {
			return err
		}
		t.OverallRecord.Wins = wins
	case 8:
		losses, err := strconv.Atoi(valStr)
		if err != nil {
			return err
		}
		t.OverallRecord.Losses = losses
	case 12:
		wins, err := strconv.Atoi(valStr)
		if err != nil {
			return err
		}
		t.ConferenceRecord.Wins = wins
	case 13:
		losses, err := strconv.Atoi(valStr)
		if err != nil {
			return err
		}
		t.ConferenceRecord.Losses = losses
	case 17:
		if valStr == "" {
			return nil
		}
		wins, err := strconv.Atoi(valStr)
		if err != nil {
			return err
		}
		if t.DivisionRecord == nil {
			t.DivisionRecord = &models.Record{}
		}
		t.DivisionRecord.Wins = wins
	case 18:
		if valStr == "" {
			return nil
		}
		losses, err := strconv.Atoi(valStr)
		if err != nil {
			return err
		}
		if t.DivisionRecord == nil {
			t.DivisionRecord = &models.Record{}
		}
		t.DivisionRecord.Losses = losses
	}
	return nil
}

func rowToTeamStanding(row []interface{}) (models.TeamStanding, error) {
	teamStanding := models.TeamStanding{}

	for i, val := range row {
		err := setTeamStandingValBasedOnColumn(&teamStanding, i, val)
		if err != nil {
			return teamStanding, err
		}
	}

	return teamStanding, nil
}

func assertToString(val interface{}) (string, error) {
	if valStr, ok := val.(string); ok {
		return valStr, nil
	}
	return "", fmt.Errorf("can't convert to string: %v", val)
}
