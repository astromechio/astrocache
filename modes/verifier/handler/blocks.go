package handler

import (
	"net/http"

	"github.com/astromechio/astrocache/execute"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/modes/common"
	"github.com/astromechio/astrocache/transport"
	"github.com/pkg/errors"
)

// ProposeAddBlockHandler adds a proposed new block and responds with the proposed prevBlock
func ProposeAddBlockHandler(chain *blockchain.Chain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config := common.ConfigFromReqCtx(r)

		proposeReq := &requests.ProposeBlockRequest{}
		proposeReq.FromRequest(r)

		if err := requests.VerifyRequest(proposeReq); err != nil {
			logger.LogError(errors.Wrap(err, "ProposeAddBlockHandler failed to VerifyRequest"))
			transport.BadRequest(w)
			return
		}

		res, err := execute.AddPendingBlock(chain, config.KeySet, proposeReq)
		if err != nil {
			logger.LogWarn(errors.Wrap(err, "ProposeAddBlock failed to AddPendingBLock"))
			if res != nil {
				transport.ReplyWithConflictJSON(w, res)
			} else {
				transport.InternalServerError(w)
			}

			return
		}

		transport.ReplyWithJSON(w, res)
	}
}
