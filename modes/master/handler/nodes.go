package handler

import (
	"net/http"

	"github.com/astromechio/astrocache/model/request"
	"github.com/astromechio/astrocache/transport"
)

// AddVerifierNodeHandler handles POST /v1/master/nodes/verifier
func AddVerifierNodeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newNodeRequest := &request.NewNodeRequest{}
		if err := newNodeRequest.FromRequest(r); err != nil {
			transport.BadRequest(w)
			return
		}

		if err := request.VerifyRequest(newNodeRequest); err != nil {
			transport.BadRequest(w)
			return
		}

		// TODO: get verifier to mine block

	}
}
