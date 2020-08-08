package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mellena1/RSC-Spreadsheet-API/data/sheets"
	"github.com/mellena1/RSC-Spreadsheet-API/handler"
	log "github.com/sirupsen/logrus"
)

// RouterCreator a handler that can return a gorilla mux router for path prefixes
type RouterCreator interface {
	Router() *mux.Router
}

func main() {
	// runHTTPServer()

	teamStandings, err := sheets.NewTeamStandingsSheet(
		context.TODO(),
		"1l99BZtpFdVB8M6xB7VJii4aAj5O33u6HUvZLGfwHB0k",
		"All Teams Data",
		os.Getenv("RSC_SHEETS_API_TOKEN"),
	)

	if err != nil {
		log.Fatalf("Error making TeamStandingsSheet: %v\n", err)
	}

	teams, err := teamStandings.GetTeamStandingsFromSheet()
	if err != nil {
		log.Fatalf("Error getting Teams: %v\n", err)
	}

	for _, t := range teams {
		fmt.Println(t)
	}
}

func runHTTPServer() {
	router := mux.NewRouter()

	childRouters := getChildRouters()
	for _, c := range childRouters {
		router.PathPrefix(c.PathPrefix).Handler(c.Child.Router())
	}

	http.Handle("/", router)
}

// ChildRouter holds a child handler that can be used for path prefixes
type ChildRouter struct {
	PathPrefix string
	Child      RouterCreator
}

func getChildRouters() []ChildRouter {
	return []ChildRouter{
		{
			PathPrefix: "/standings/",
			Child:      &handler.StandingsHandler{},
		},
	}
}
