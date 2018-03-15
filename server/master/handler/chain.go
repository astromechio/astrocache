package handler

import (
	"net/http"

	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/pkg/errors"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/transport"
	"github.com/gorilla/mux"
)

// GetEntireChainHandler returns the entire chain for a node to verify and store
func GetEntireChainHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transport.ReplyWithJSON(w, app.Chain.Blocks)
	}
}

const afterKey = "after"

// GetBlocksAfterHandler handles blocks after ID requests
func GetBlocksAfterHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		afterID := mux.Vars(r)[afterKey]

		blocks := app.Chain.BlocksAfterID(afterID)
		if blocks == nil {
			transport.InternalServerError(w)
			return
		}

		transport.ReplyWithJSON(w, blocks)
	}
}

// ReserveIDHandler handles blocks after ID requests
func ReserveIDHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reserveReq := &requests.ReserveIDRequest{}
		reserveReq.FromRequest(r)

		if err := requests.VerifyRequest(reserveReq); err != nil {
			logger.LogError(errors.Wrap(err, "ReserveIDHandler failed to VerifyRequest"))
			transport.BadRequest(w)
			return
		}

		errChan, reserveJob := app.Chain.ReserveBlockID(reserveReq.ProposingNID)
		if err := <-errChan; err != nil {
			logger.LogError(errors.New("ReserveIDHandler failed to ReserveBlockID, got empty blockID"))
			transport.Conflict(w)
			return
		}

		resp := &requests.ReserveIDResponse{
			BlockID: reserveJob.BlockID,
		}

		transport.ReplyWithJSON(w, resp)
	}
}
