package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/astronomy/ap-manager/helpers"
	"gitlab.com/astronomy/ap-manager/storage"
)

type SessionAddRequest struct {
	FileName string `json:"filename"`
}

func SessionAddRouteHandler(e *Env, w http.ResponseWriter, r *http.Request) error {

	var sr SessionAddRequest

	err := json.NewDecoder(r.Body).Decode(&sr)
	if err != nil {
		return GetStatusErrorForCode(http.StatusBadRequest)
	}

	s := helpers.ProcessLogFile(sr.FileName)
	resp := storage.CreateSessions(e.DB, s)
	response, errr := json.Marshal(resp)
	return HandlerParseResponse(w, response, errr)
}
