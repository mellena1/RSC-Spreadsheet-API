package handler

import (
	"encoding/json"
	"net/http"
)

// errorResp a model to respond to users with for errors
type errorResp struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, errorMsg string, statuscode int) {
	msg, _ := json.Marshal(&errorResp{Error: errorMsg})
	w.WriteHeader(statuscode)
	w.Write(msg)
}
