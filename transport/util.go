package transport

import (
	"encoding/json"
	"net/http"

	"github.com/astromechio/astrocache/logger"
	"github.com/pkg/errors"
)

// ReplyWithJSON replies with json and 200
func ReplyWithJSON(w http.ResponseWriter, value interface{}) {
	response, err := json.Marshal(value)
	if err != nil {
		logger.LogError(errors.Wrap(err, "ReplyWithJSON failed to Marshal"))
		InternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ReplyWithConflictJSON replies with json and 200
func ReplyWithConflictJSON(w http.ResponseWriter, value interface{}) {
	response, err := json.Marshal(value)
	if err != nil {
		logger.LogError(errors.Wrap(err, "ReplyWithJSON failed to Marshal"))
		InternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	w.Write(response)
}

// InternalServerError respoonds with 500
func InternalServerError(w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// NotFound responds with 404
func NotFound(w http.ResponseWriter) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

// BadRequest responds with 400
func BadRequest(w http.ResponseWriter) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
}

// Forbidden responds with 403
func Forbidden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}
