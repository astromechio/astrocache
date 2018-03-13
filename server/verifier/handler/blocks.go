package handler

import (
	"fmt"
	"net/http"

	"github.com/astromechio/astrocache/model/blockchain"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/transport"
	"github.com/pkg/errors"
)

// ProposeAddBlockHandler adds a proposed new block
func ProposeAddBlockHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chain := app.Chain

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
func CheckBlockHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chain := app.Chain

		checkReq := &requests.CheckBlockRequest{}
		checkReq.FromRequest(r)

		if err := requests.VerifyRequest(checkReq); err != nil {
			logger.LogError(errors.Wrap(err, "ProposeAddBlockHandler failed to VerifyRequest"))
			transport.BadRequest(w)
			return
		}

		if checkLastAndProposed(chain, checkReq.Block) {
			transport.Ok(w)
			return
		}

		var block *blockchain.Block

		proposedChan := chain.GetNextProposed() // blocks until proposed is set
		committedChan := chain.GetNextCommitted()

		for true {
			select {
			case block = <-proposedChan:
				break
			case block = <-committedChan:
				break
			default:
				if checkLastAndProposed(chain, checkReq.Block) {
					logger.LogInfo("Checked against chain; succeeded")
					transport.Ok(w)
					return
				}
			}

		}

		logger.LogInfo("Checking against new block with ID " + block.ID)

		if !block.IsSameAsBlock(checkReq.Block) {
			logger.LogError(fmt.Errorf("CheckBlockHandler failed to check block %s, is not same as new block", checkReq.Block.ID))
			transport.Conflict(w)
			return
		}

		transport.Ok(w)
	}
}

func checkLastAndProposed(chain *blockchain.Chain, block *blockchain.Block) bool {
	lastBlock := chain.LastBlock()
	if lastBlock.IsSameAsBlock(block) {
		return true
	}

	proposed := chain.Proposed
	if proposed != nil && proposed.IsSameAsBlock(block) {
		return true
	}

	return false
}
