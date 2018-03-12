package handler

import (
	"net/http"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/transport"
)

// GetEntireChainHandler returns the entire chain for a node to verify and store
func GetEntireChainHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transport.ReplyWithJSON(w, app.Chain.Blocks)
	}
}
