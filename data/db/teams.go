package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mellena1/RSC-Spreadsheet-API/data/models"
	log "github.com/sirupsen/logrus"
)

// ErrInvalidTypeForQuery is returned if a given query is the wrong type
var ErrInvalidTypeForQuery error = errors.New("Invalid query param type")

func createWhereQuery(startingNum int, fieldName string, lenVals int, separator string) string {
	queryStr := ""

	for i := 0; i < lenVals; i++ {
		queryStr += fmt.Sprintf("%s=$%d %s ", fieldName, startingNum, separator)
		startingNum++
	}

	return strings.TrimSuffix(queryStr, fmt.Sprintf(" %s ", separator))
}

func createWhereQueryWithAnds(startingNum int, fieldName string, lenVals int) string {
	return createWhereQuery(startingNum, fieldName, lenVals, "AND")
}

func createWhereQueryWithOrs(startingNum int, fieldName string, lenVals int) string {
	return createWhereQuery(startingNum, fieldName, lenVals, "OR")
}

func stringSliceToInterfaceSlice(vals []string) []interface{} {
	interfaceVals := make([]interface{}, len(vals))
	for i, v := range vals {
		interfaceVals[i] = interface{}(v)
	}
	return interfaceVals
}

func stringIsInt(s string) bool {
	if _, err := strconv.Atoi(s); err != nil {
		return false
	}
	return true
}

type GetAllTeamsQuery struct {
	TeamIDs     []string
	Names       []string
	Franchises  []string
	Conferences []string
	Tiers       []string
	Divisions   []string
}

func (q GetAllTeamsQuery) buildQueryStr(startingNum int) (string, []interface{}, error) {
	queryStr := ""
	params := []string{}

	if len(q.TeamIDs) > 0 {
		for _, id := range q.TeamIDs {
			if !stringIsInt(id) {
				return "", nil, ErrInvalidTypeForQuery
			}
		}
		queryStr += fmt.Sprintf("(%s)", createWhereQueryWithOrs(startingNum, "team_id", len(q.TeamIDs)))
		params = append(params, q.TeamIDs...)
		startingNum += len(q.TeamIDs)
	}
	if len(q.Names) > 0 {
		if queryStr != "" {
			queryStr += " AND "
		}
		queryStr += fmt.Sprintf("(%s)", createWhereQueryWithOrs(startingNum, "name", len(q.Names)))
		params = append(params, q.Names...)
		startingNum += len(q.Names)
	}
	if len(q.Franchises) > 0 {
		if queryStr != "" {
			queryStr += " AND "
		}
		queryStr += fmt.Sprintf("(%s)", createWhereQueryWithOrs(startingNum, "franchise", len(q.Franchises)))
		params = append(params, q.Franchises...)
		startingNum += len(q.Franchises)
	}
	if len(q.Conferences) > 0 {
		if queryStr != "" {
			queryStr += " AND "
		}
		queryStr += fmt.Sprintf("(%s)", createWhereQueryWithOrs(startingNum, "conference", len(q.Conferences)))
		params = append(params, q.Conferences...)
		startingNum += len(q.Conferences)
	}
	if len(q.Tiers) > 0 {
		if queryStr != "" {
			queryStr += " AND "
		}
		queryStr += fmt.Sprintf("(%s)", createWhereQueryWithOrs(startingNum, "tier", len(q.Tiers)))
		params = append(params, q.Tiers...)
		startingNum += len(q.Tiers)
	}
	if len(q.Divisions) > 0 {
		if queryStr != "" {
			queryStr += " AND "
		}
		queryStr += fmt.Sprintf("(%s)", createWhereQueryWithOrs(startingNum, "division", len(q.Divisions)))
		params = append(params, q.Divisions...)
		startingNum += len(q.Divisions)
	}

	if queryStr != "" {
		queryStr = "WHERE " + queryStr
	}

	return queryStr, stringSliceToInterfaceSlice(params), nil
}

func (db *DB) GetAllTeams(query GetAllTeamsQuery) ([]models.Team, error) {
	conditionalStr, params, err := query.buildQueryStr(1)
	if err != nil {
		log.Warnf("Error making sql query from GetAllTeamsQuery %+v", query)
		return nil, err
	}

	sqlQuery := fmt.Sprintf(`
		SELECT team_id, name, franchise, conference, tier, division FROM team %s;
	`, conditionalStr)

	rows, err := db.sqlDB.Query(sqlQuery, params...)
	if err != nil {
		log.Errorf("Error getting all teams from db: %v", err)
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
