package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_createWhereQueryWithAnds(t *testing.T) {
	tests := []struct {
		name        string
		startingNum int
		fieldName   string
		lenVals     int
		expected    string
	}{
		{
			name:        "Standard call",
			startingNum: 1,
			fieldName:   "myfield",
			lenVals:     5,
			expected:    "myfield=$1 AND myfield=$2 AND myfield=$3 AND myfield=$4 AND myfield=$5",
		},
		{
			name:        "Empty call",
			startingNum: 1,
			fieldName:   "myfield",
			lenVals:     0,
			expected:    "",
		},
	}

	for _, test := range tests {
		actual := createWhereQueryWithAnds(test.startingNum, test.fieldName, test.lenVals)
		require.Equalf(t, test.expected, actual, "test %q failed", test.name)
	}
}

func Test_stringSliceToInterfaceSlice(t *testing.T) {
	expected := []interface{}{"a", "b", "c"}
	actual := stringSliceToInterfaceSlice([]string{"a", "b", "c"})

	require.Equal(t, expected, actual)
}

func Test_GetAllTeamsQuery_buildQueryStr(t *testing.T) {
	tests := []struct {
		name           string
		startingNum    int
		query          GetAllTeamsQuery
		expectedStr    string
		expectedParams []interface{}
		expectedErr    error
	}{
		{
			name:        "All fields",
			startingNum: 1,
			query: GetAllTeamsQuery{
				TeamIDs:     []string{"1", "2"},
				Names:       []string{"Care Bears"},
				Franchises:  []string{"Bear Den"},
				Conferences: []string{"Solar", "Lunar"},
				Tiers:       []string{"Master", "Elite"},
				Divisions:   []string{"Solar Mountain"},
			},
			expectedStr:    "WHERE (team_id=$1 OR team_id=$2) AND (name=$3) AND (franchise=$4) AND (conference=$5 OR conference=$6) AND (tier=$7 OR tier=$8) AND (division=$9)",
			expectedParams: []interface{}{"1", "2", "Care Bears", "Bear Den", "Solar", "Lunar", "Master", "Elite", "Solar Mountain"},
		},
		{
			name:        "Some fields",
			startingNum: 1,
			query: GetAllTeamsQuery{
				Conferences: []string{"Solar"},
				Divisions:   []string{"Solar Mountain"},
			},
			expectedStr:    "WHERE (conference=$1) AND (division=$2)",
			expectedParams: []interface{}{"Solar", "Solar Mountain"},
		},
		{
			name:           "Empty",
			startingNum:    1,
			query:          GetAllTeamsQuery{},
			expectedStr:    "",
			expectedParams: []interface{}{},
		},
		{
			name:        "Invalid team id",
			startingNum: 1,
			query: GetAllTeamsQuery{
				TeamIDs: []string{"abc"},
			},
			expectedErr: ErrInvalidTypeForQuery,
		},
	}

	for _, test := range tests {
		actualStr, actualParams, actualErr := test.query.buildQueryStr(test.startingNum)
		require.Equalf(t, test.expectedStr, actualStr, "test %q failed", test.name)
		require.Equalf(t, test.expectedParams, actualParams, "test %q failed", test.name)
		require.Equalf(t, test.expectedErr, actualErr, "test %q failed", test.name)
	}
}
