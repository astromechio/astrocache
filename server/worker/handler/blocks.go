package handler

import (
	"net/http"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/transport"
	"github.com/pkg/errors"
)

// AddBlockHandler handles adding new blocks
func AddBlockHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chain := app.Chain

		proposeReq := &requests.ProposeBlockRequest{}
		proposeReq.FromRequest(r)

		if err := requests.VerifyRequest(proposeReq); err != nil {
			logger.LogError(errors.Wrap(err, "ProposeAddBlockHandler failed to VerifyRequest"))
			transport.BadRequest(w)
			return
		}

		errChan := chain.VerifyBlockUnchecked(proposeReq.Block)
		if err := <-errChan; err != nil {
			logger.LogError(errors.Wrap(err, "ProposeAddBlockHandler failed to AddNewBlock"))
			transport.Conflict(w)
			return
		}

		transport.Ok(w)
	}
}
