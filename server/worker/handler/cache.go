package handler

import (
	"net/http"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/send"
	"github.com/astromechio/astrocache/transport"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// SetValueHandler handles value set requests
func SetValueHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setValReq := &requests.SetValueRequest{}
		setValReq.FromRequest(r)

		if err := setValReq.Verify(); err != nil {
			logger.LogError(errors.Wrap(err, "SetValueHandler failed to Verify"))
			transport.BadRequest(w)
			return
		}

		if err := send.SetValue(setValReq, app.NodeList.RandomVerifier()); err != nil {
			logger.LogError(errors.Wrap(err, "SetValueHandler failed to SetValue"))
			transport.InternalServerError(w)
			return
		}

		transport.Ok(w)
	}
}

// GetValueHandler handles value set requests
func GetValueHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)[requests.KeyRequestKey]

		val := app.Cache.ValueForKey(key)
		if val == "" {
			transport.NotFound(w)
			return
		}

		w.Write([]byte(val))
		w.WriteHeader(http.StatusOK)
	}
}
