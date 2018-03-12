package handler

import (
	"fmt"
	"net/http"

	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/server"
	"github.com/astromechio/astrocache/transport"
	"github.com/pkg/errors"
)

// ProposeAddBlockHandler adds a proposed new block
func ProposeAddBlockHandler(config *server.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chain := config.Chain

		proposeReq := &requests.ProposeBlockRequest{}
		proposeReq.FromRequest(r)

		if err := requests.VerifyRequest(proposeReq); err != nil {
			logger.LogError(errors.Wrap(err, "ProposeAddBlockHandler failed to VerifyRequest"))
			transport.BadRequest(w)
			return
		}

		errChan := chain.AddNewBlock(proposeReq.Block)
		if err := <-errChan; err != nil {
			logger.LogError(errors.Wrap(err, "ProposeAddBlockHandler failed to AddNewBlock"))
			transport.Conflict(w)
		}

		transport.Ok(w)
	}
}

// CheckBlockHandler adds a proposed new block and responds with the proposed prevBlock
func CheckBlockHandler(config *server.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chain := config.Chain

		checkReq := &requests.CheckBlockRequest{}
		checkReq.FromRequest(r)

		if err := requests.VerifyRequest(checkReq); err != nil {
			logger.LogError(errors.Wrap(err, "ProposeAddBlockHandler failed to VerifyRequest"))
			transport.BadRequest(w)
			return
		}

		lastBlock := chain.LastBlock()
		if !lastBlock.IsSameAsBlock(checkReq.Block) {
			proposed := chain.Proposed
			if !proposed.IsSameAsBlock(checkReq.Block) {
				logger.LogError(fmt.Errorf("CheckBlockHandler failed to check block %s, is not same as LastBlock or Proposed", checkReq.Block.ID))
				transport.Conflict(w)
				return
			}
		}

		transport.Ok(w)
	}
}
