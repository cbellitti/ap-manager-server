package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func HealthCheckRouteHandler(e *Env, w http.ResponseWriter, r *http.Request) error {
	log.Println("Health check handler")

	response, errr := json.Marshal(HealthCheckResponse{"SUCCESS - UF server running on port 8085"})

	return HandlerParseResponse(w, response, errr)

}
