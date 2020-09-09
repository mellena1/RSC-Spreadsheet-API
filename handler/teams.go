package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mellena1/RSC-Spreadsheet-API/data/db"
	"github.com/mellena1/RSC-Spreadsheet-API/data/models"
	log "github.com/sirupsen/logrus"
)

// TeamHandler has all routes for team related queries
type TeamHandler struct {
	DB db.Datastore
}

// AddRoutes adds all of it's routes to the router
func (t *TeamHandler) AddRoutes(router *mux.Router) {
	if t.DB == nil {
		log.Fatal("TeamHandler.DB is nil!")
	}

	router.HandleFunc("", t.getAllTeams).Methods("GET")
	router.HandleFunc("/", t.getAllTeams).Methods("GET")
	router.HandleFunc("/{id}", t.getTeam).Methods("GET")
}

// teamsListResp A List of Teams
// swagger:response teamsListResp
type teamsListResp struct {
	// in: body
	Teams []models.Team `json:"teams"`
}

// swagger:parameters getAllTeams
type getAllTeamsParams struct {
	ID         []int    `json:"id"`
	Name       []string `json:"name"`
	Franchise  []string `json:"franchise"`
	Conference []string `json:"conference"`
	Tier       []string `json:"tier"`
	Division   []string `json:"division"`
}

// getAllTeams is a handler to get all teams
// swagger:route GET /team getAllTeams
// Gets all RSC teams
// Responses:
// 	 200: teamsListResp
func (t *TeamHandler) getAllTeams(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Errorf("Invalid URL query string: %s", err)
		writeError(w, "Invalid query", http.StatusBadRequest)
		return
	}

	query := db.GetAllTeamsQuery{
		TeamIDs:     r.Form["id"],
		Names:       r.Form["name"],
		Franchises:  r.Form["franchise"],
		Conferences: r.Form["conference"],
		Tiers:       r.Form["tier"],
		Divisions:   r.Form["division"],
	}

	teams, err := t.DB.GetAllTeams(query)
	if err == db.ErrInvalidTypeForQuery {
		log.Warn("Invalid query param for team")
		writeError(w, "Team IDs must be integers", http.StatusBadRequest)
		return
	} else if err != nil {
		log.Errorf("Unable to fetch teams from db: %s", err)
		writeError(w, "Failed to fetch teams from db", http.StatusInternalServerError)
		return
	}

	msg, err := json.Marshal(&teamsListResp{Teams: teams})
	if err != nil {
		log.Errorf("Unable to marshal teams: %s", err)
		writeError(w, "Error sending team", http.StatusInternalServerError)
		return
	}
	w.Write(msg)
}

// getTeam is a handler to get a single team
// swagger:route GET /team/{id} getTeam
// Gets a RSC team by id
// Responses:
// 	 200: models_Team
func (t *TeamHandler) getTeam(w http.ResponseWriter, r *http.Request) {
	teamID := mux.Vars(r)["id"]
	query := db.GetAllTeamsQuery{
		TeamIDs: []string{teamID},
	}

	teams, err := t.DB.GetAllTeams(query)
	if err == db.ErrInvalidTypeForQuery {
		log.Warnf("Invalid team id: %s", teamID)
		writeError(w, "Team ID must be an integer", http.StatusBadRequest)
		return
	} else if err != nil {
		log.Errorf("Unable to fetch team from db: %s", err)
		writeError(w, "Failed to fetch team from db", http.StatusInternalServerError)
		return
	}

	if len(teams) == 0 {
		writeError(w, "Team not found", http.StatusNotFound)
		return
	}

	msg, err := json.Marshal(&teams[0])
	if err != nil {
		log.Errorf("Unable to marshal teams: %s", err)
		writeError(w, "Error sending team", http.StatusInternalServerError)
		return
	}
	w.Write(msg)
}
