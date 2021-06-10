package alcon

import (
	"net/http"

	httprouter "github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// defaultHandler default endpoint handler
func DefaultHandler(w http.ResponseWriter, r *http.Request, router httprouter.Params) {
	log.WithFields(log.Fields{
		"url path": r.URL.Path,
	}).Info("Default URL path")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}
