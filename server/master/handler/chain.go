package handler

import (
	"net/http"

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
