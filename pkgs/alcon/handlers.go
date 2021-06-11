package alcon

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// DefaultHandler default endpoint handler
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{
		"url path": r.URL.Path,
	}).Info("GET for Default URL path")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode([]string{"Alive"})

	return
}
