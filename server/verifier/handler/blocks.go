package handler

import (
	"fmt"
	"net/http"

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

		errChan := chain.VerifyProposedBlock(proposeReq.Block, proposeReq.ProposingNID)
		if err := <-errChan; err != nil {
			logger.LogError(errors.Wrap(err, "ProposeAddBlockHandler failed to AddNewBlock"))
			transport.Conflict(w)
			return
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

		if !chain.HasProposedOrCommittedBlock(checkReq.Block) {
			logger.LogError(fmt.Errorf("CheckBlockHandler failed to HasProposedOrCommittedBlock for block with ID %q", checkReq.Block.ID))
			transport.Conflict(w)
			return
		}

		transport.Ok(w)
	}
}

// // If that failed, we now have to do a convoluted channel dance to wait for the right block to come around
// var block *blockchain.Block

// proposedChan := chain.GetNextProposed() // blocks until proposed is set
// committedChan := chain.GetNextCommitted()

// for true {
// 	select {
// 	case block = <-proposedChan:
// 		break
// 	case block = <-committedChan:
// 		break
// 	default:
// 		if checkLastAndProposed(chain, checkReq.Block) {
// 			logger.LogInfo("Checked against chain; succeeded")
// 			transport.Ok(w)
// 			return
// 		}
// 	}

// }

// logger.LogInfo("Checking against new block with ID " + block.ID)

// if !block.IsSameAsBlock(checkReq.Block) {
// 	logger.LogError(fmt.Errorf("CheckBlockHandler failed to check block %q, is not same as new block", checkReq.Block.ID))
// 	transport.Conflict(w)
// 	return
// }
