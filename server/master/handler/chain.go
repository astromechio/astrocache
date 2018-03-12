package handler

import (
	"net/http"

	"github.com/astromechio/astrocache/server"
	"github.com/astromechio/astrocache/transport"
)

// GetEntireChainHandler returns the entire chain for a node to verify and store
func GetEntireChainHandler(config *server.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transport.ReplyWithJSON(w, config.Chain.Blocks)
	}
}
