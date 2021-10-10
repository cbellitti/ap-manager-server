package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/astronomy/ap-manager/storage"
)

func SessionsRouteHandler(e *Env, w http.ResponseWriter, r *http.Request) error {
	s := storage.ReadSessions(e.DB)
	response, errr := json.Marshal(s)
	return HandlerParseResponse(w, response, errr)
}
