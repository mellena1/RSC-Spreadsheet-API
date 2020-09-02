package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mellena1/RSC-Spreadsheet-API/data/db"
	"github.com/mellena1/RSC-Spreadsheet-API/data/sheets"
	"github.com/mellena1/RSC-Spreadsheet-API/handler"
	log "github.com/sirupsen/logrus"
)

// RouterCreator a handler that can return a gorilla mux router for path prefixes
type RouterCreator interface {
	AddRoutes(*mux.Router)
}

func fatalIfMissingEnvVar(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	log.Fatalf("Must set env var: %s", key)
	return ""
}

func getEnvOrDefault(key, _default string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return _default
}

func main() {
	mydb := makeDB()
	defer mydb.Close()

	router := makeHTTPRouter(mydb)
	log.Info("Serving on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	if err := mydb.Close(); err != nil {
		log.Fatalf("Error closing db: %v\n", err)
	}
}

func makeDB() *db.DB {
	teamStandings, err := sheets.NewTeamStandingsSheet(
		context.TODO(),
		"1l99BZtpFdVB8M6xB7VJii4aAj5O33u6HUvZLGfwHB0k",
		"All Teams Data",
		fatalIfMissingEnvVar("RSC_SHEETS_API_TOKEN"),
	)
	if err != nil {
		log.Fatalf("Error making TeamStandingsSheet: %v\n", err)
	}

	dbStr := fmt.Sprintf(
		"postgres://%s:%s@%s?sslmode=disable",
		getEnvOrDefault("DB_USER", "postgres"),
		getEnvOrDefault("DB_PASS", "password"),
		fatalIfMissingEnvVar("DB_HOST"),
	)
	mydb, err := db.NewDB(dbStr, teamStandings)
	if err != nil {
		log.Fatalf("Error making db: %v\n", err)
	}

	return mydb
}

func makeHTTPRouter(_db *db.DB) http.Handler {
	router := mux.NewRouter()

	childRouters := getChildRouters(_db)
	for _, c := range childRouters {
		subR := router.PathPrefix(c.PathPrefix).Subrouter()
		c.Child.AddRoutes(subR)
	}

	loggingRouterHandler := gorillaHandlers.LoggingHandler(os.Stdout, router)

	return loggingRouterHandler
}

// ChildRouter holds a child handler that can be used for path prefixes
type ChildRouter struct {
	PathPrefix string
	Child      RouterCreator
}

func getChildRouters(_db *db.DB) []ChildRouter {
	return []ChildRouter{
		{
			PathPrefix: "/team",
			Child: &handler.TeamHandler{
				DB: _db,
			},
		},
	}
}
