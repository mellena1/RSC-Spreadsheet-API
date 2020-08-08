package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// StandingsHandler has all routes for standings related things
type StandingsHandler struct {
}

// Router returns a gorilla router
func (s *StandingsHandler) Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", s.getAllTeamsRecords)

	return router
}

func (s *StandingsHandler) getAllTeamsRecords(w http.ResponseWriter, r *http.Request) {

}
