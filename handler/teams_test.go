package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mellena1/RSC-Spreadsheet-API/data/db"
	"github.com/mellena1/RSC-Spreadsheet-API/data/models"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type datastoreEmptyMock struct {
	db.Datastore
}

func makeReq(method, host string) *http.Request {
	req, _ := http.NewRequest(method, host, nil)
	return req
}

func Test_AddRoutes(t *testing.T) {
	router := mux.NewRouter()

	tHandler := TeamHandler{DB: datastoreEmptyMock{}}
	tHandler.AddRoutes(router)

	tests := []struct {
		req          *http.Request
		expected     http.HandlerFunc
		expectedVars map[string]string
	}{
		{
			req:      makeReq("GET", ""),
			expected: tHandler.getAllTeams,
		},
		{
			req:      makeReq("GET", "/"),
			expected: tHandler.getAllTeams,
		},
		{
			req:          makeReq("GET", "/10"),
			expected:     tHandler.getTeam,
			expectedVars: map[string]string{"id": "10"},
		},
	}
	for _, test := range tests {
		routeMatch := &mux.RouteMatch{}
		matched := router.Match(test.req, routeMatch)
		require.Equal(t, true, matched)
		// use sprintf to compare function addresses
		require.Equal(t, fmt.Sprintf("%v", test.expected), fmt.Sprintf("%v", routeMatch.Handler))
		if test.expectedVars != nil {
			require.Equal(t, test.expectedVars, routeMatch.Vars)
		}
	}
}

func Test_AddRoutes_NilDB(t *testing.T) {
	origExitFunc := log.StandardLogger().ExitFunc
	defer func() { log.StandardLogger().ExitFunc = origExitFunc }()
	var fatal bool
	log.StandardLogger().ExitFunc = func(int) { fatal = true }

	tHandler := TeamHandler{}
	tHandler.AddRoutes(mux.NewRouter())

	require.Equal(t, true, fatal)
}

type getAllTeamsMockDB struct {
	db.Datastore

	t                *testing.T
	expectedQueryVal db.GetAllTeamsQuery

	resp []models.Team
	err  error
}

func (d getAllTeamsMockDB) GetAllTeams(query db.GetAllTeamsQuery) ([]models.Team, error) {
	require.Equal(d.t, d.expectedQueryVal, query)

	return d.resp, d.err
}

func strPointer(s string) *string {
	return &s
}

func Test_getAllTeams(t *testing.T) {
	tests := []struct {
		mockDB             db.Datastore
		requestPath        string
		requestMethod      string
		expectedResp       string
		expectedStatusCode int
	}{
		{
			requestPath:        "/?id=1&id=2&name=A&franchise=B&conference=C&tier=D&division=E",
			requestMethod:      "GET",
			expectedResp:       `{"teams":[{"id":"1","name":"A","franchise":"B","tier":"D","conference":"C","division":"E"},{"id":"2","name":"A","franchise":"B","tier":"D","conference":"C","division":"E"}]}`,
			expectedStatusCode: 200,
			mockDB: getAllTeamsMockDB{
				t: t,
				expectedQueryVal: db.GetAllTeamsQuery{
					TeamIDs:     []string{"1", "2"},
					Names:       []string{"A"},
					Franchises:  []string{"B"},
					Conferences: []string{"C"},
					Tiers:       []string{"D"},
					Divisions:   []string{"E"},
				},
				resp: []models.Team{
					{TeamID: "1", Name: "A", Franchise: "B", Conference: "C", Tier: "D", Division: strPointer("E")},
					{TeamID: "2", Name: "A", Franchise: "B", Conference: "C", Tier: "D", Division: strPointer("E")},
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		tHandler := TeamHandler{DB: test.mockDB}
		router := mux.NewRouter()
		tHandler.AddRoutes(router)
		server := httptest.NewServer(router)
		t.Cleanup(server.Close)

		url := fmt.Sprintf("%s%s", server.URL, test.requestPath)
		req, _ := http.NewRequest(test.requestMethod, url, nil)
		actual, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { actual.Body.Close() })
		require.Equal(t, test.expectedStatusCode, actual.StatusCode)
		body, err := ioutil.ReadAll(actual.Body)
		require.NoError(t, err)
		require.Equal(t, test.expectedResp, string(body))
	}
}
