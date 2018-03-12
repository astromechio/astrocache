package handler

import (
	"net/http"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/transport"
	"github.com/pkg/errors"
)

// SetValueHandler handles value set requests
func SetValueHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setValReq := requests.SetValueRequest{}
		setValReq.FromRequest(r)

		if err := setValReq.Verify(); err != nil {
			logger.LogError(errors.Wrap(err, "SetValueHandler failed to Verify"))
			transport.BadRequest(w)
			return
		}

	}
}
