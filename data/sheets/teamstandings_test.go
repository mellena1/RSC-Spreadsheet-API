package sheets

import (
	"errors"
	"testing"

	"github.com/mellena1/RSC-Spreadsheet-API/data/models"
	"github.com/stretchr/testify/require"
)

func stringPtr(s string) *string {
	return &s
}

func Test_setTeamStandingValBasedOnColumn(t *testing.T) {
	tests := []struct {
		name        string
		colIndex    int
		val         interface{}
		expectedT   *TeamStanding
		expectedErr error
	}{
		{
			name:        "not a string",
			colIndex:    0,
			val:         interface{}(1),
			expectedT:   &TeamStanding{},
			expectedErr: errors.New("can't convert to string: 1"),
		},
		{
			name:        "Tier",
			colIndex:    0,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{Team: models.Team{Tier: "1"}},
			expectedErr: nil,
		},
		{
			name:        "Franchise",
			colIndex:    1,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{Team: models.Team{Franchise: "1"}},
			expectedErr: nil,
		},
		{
			name:        "Name",
			colIndex:    2,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{Team: models.Team{Name: "1"}},
			expectedErr: nil,
		},
		{
			name:        "Conference",
			colIndex:    3,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{Team: models.Team{Conference: "1"}},
			expectedErr: nil,
		},
		{
			name:        "Division",
			colIndex:    4,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{Team: models.Team{Division: stringPtr("1")}},
			expectedErr: nil,
		},
		{
			name:        "Division N/A",
			colIndex:    4,
			val:         interface{}("N/A"),
			expectedT:   &TeamStanding{Team: models.Team{Division: nil}},
			expectedErr: nil,
		},
		{
			name:        "Overall.Wins",
			colIndex:    7,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{OverallRecord: Record{Wins: 1}},
			expectedErr: nil,
		},
		{
			name:        "Overall.Wins not a number",
			colIndex:    7,
			val:         interface{}("abc"),
			expectedT:   &TeamStanding{},
			expectedErr: errors.New(`strconv.Atoi: parsing "abc": invalid syntax`),
		},
		{
			name:        "Overall.Losses",
			colIndex:    8,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{OverallRecord: Record{Losses: 1}},
			expectedErr: nil,
		},
		{
			name:        "Overall.Losses not a number",
			colIndex:    8,
			val:         interface{}("abc"),
			expectedT:   &TeamStanding{},
			expectedErr: errors.New(`strconv.Atoi: parsing "abc": invalid syntax`),
		},
		{
			name:        "ConferenceRecord.Wins",
			colIndex:    12,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{ConferenceRecord: Record{Wins: 1}},
			expectedErr: nil,
		},
		{
			name:        "ConferenceRecord.Wins not a number",
			colIndex:    12,
			val:         interface{}("abc"),
			expectedT:   &TeamStanding{},
			expectedErr: errors.New(`strconv.Atoi: parsing "abc": invalid syntax`),
		},
		{
			name:        "ConferenceRecord.Losses",
			colIndex:    13,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{ConferenceRecord: Record{Losses: 1}},
			expectedErr: nil,
		},
		{
			name:        "ConferenceRecord.Losses not a number",
			colIndex:    13,
			val:         interface{}("abc"),
			expectedT:   &TeamStanding{},
			expectedErr: errors.New(`strconv.Atoi: parsing "abc": invalid syntax`),
		},
		{
			name:        "DivisionRecord.Wins",
			colIndex:    17,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{DivisionRecord: &Record{Wins: 1}},
			expectedErr: nil,
		},
		{
			name:        "DivisionRecord.Wins empty str",
			colIndex:    17,
			val:         interface{}(""),
			expectedT:   &TeamStanding{},
			expectedErr: nil,
		},
		{
			name:        "DivisionRecord.Wins not a number",
			colIndex:    17,
			val:         interface{}("abc"),
			expectedT:   &TeamStanding{},
			expectedErr: errors.New(`strconv.Atoi: parsing "abc": invalid syntax`),
		},
		{
			name:        "DivisionRecord.Losses",
			colIndex:    18,
			val:         interface{}("1"),
			expectedT:   &TeamStanding{DivisionRecord: &Record{Losses: 1}},
			expectedErr: nil,
		},
		{
			name:        "DivisionRecord.Losses empty str",
			colIndex:    18,
			val:         interface{}(""),
			expectedT:   &TeamStanding{},
			expectedErr: nil,
		},
		{
			name:        "DivisionRecord.Losses not a number",
			colIndex:    18,
			val:         interface{}("abc"),
			expectedT:   &TeamStanding{},
			expectedErr: errors.New(`strconv.Atoi: parsing "abc": invalid syntax`),
		},
	}

	for _, test := range tests {
		standing := &TeamStanding{}
		err := setTeamStandingValBasedOnColumn(standing, test.colIndex, test.val)
		require.Equalf(t, test.expectedT, standing, "%q wrong T val", test.name)
		if test.expectedErr == nil {
			require.NoErrorf(t, err, "%q should not error", test.name)
		} else {
			require.Equalf(t, test.expectedErr.Error(), err.Error(), "%q wrong error", test.name)
		}
	}
}

func Test_rowToTeamStanding(t *testing.T) {
	standing, _ := rowToTeamStanding([]interface{}{"tier", "franchise", "name", "conf"})
	require.Equal(t, TeamStanding{
		Team: models.Team{
			Tier:       "tier",
			Franchise:  "franchise",
			Name:       "name",
			Conference: "conf",
		},
	}, standing)

	_, err := rowToTeamStanding([]interface{}{"tier", "franchise", "name", "conf", "", "", "", "abc"})
	require.EqualError(t, err, `strconv.Atoi: parsing "abc": invalid syntax`)
}

func Test_assertToString(t *testing.T) {
	s, _ := assertToString(interface{}("abc"))
	require.Equal(t, "abc", s)

	_, err := assertToString(interface{}(1))
	require.EqualError(t, err, "can't convert to string: 1")
}
