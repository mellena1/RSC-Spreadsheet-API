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
	DB *db.DB
}

// AddRoutes adds all of it's routes to the router
func (t *TeamHandler) AddRoutes(router *mux.Router) {
	if t.DB == nil {
		log.Fatal("TeamHandler.DB is nil!")
	}

	router.HandleFunc("", t.getAllTeams).Methods("GET")
	router.HandleFunc("/", t.getAllTeams).Methods("GET")
}

type teamsListResp struct {
	Teams []models.Team `json:"teams"`
}

func (t *TeamHandler) getAllTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := t.DB.GetAllTeams()
	if err != nil {
		log.Errorf("Unable to fetch teams from db: %s", err)
		writeError(w, "Failed to fetch teams from db", http.StatusInternalServerError)
		return
	}

	msg, err := json.Marshal(&teamsListResp{Teams: teams})
	if err != nil {
		log.Errorf("Unable to marshal teams: %s", err)
		writeError(w, "Error sending teams", http.StatusInternalServerError)
		return
	}
	w.Write(msg)
}
