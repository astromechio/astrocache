package handler

import (
	"crypto/rand"
	"net/http"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/model/blockchain"

	"github.com/astromechio/astrocache/logger"
	"github.com/pkg/errors"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model/actions"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/transport"
)

// AddVerifierNodeHandler handles POST /v1/master/nodes/verifier
func AddVerifierNodeHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newNodeRequest := &requests.NewNodeRequest{}
		if err := newNodeRequest.FromRequest(r); err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to FromRequest"))
			transport.BadRequest(w)
			return
		}

		if err := requests.VerifyRequest(newNodeRequest); err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to VerifyRequest"))
			transport.BadRequest(w)
			return
		}

		joinCode := app.ValueForKey(config.AppJoinCodeKey)
		if newNodeRequest.JoinCode != joinCode {
			logger.LogError(errors.New("AddVerifierNodeHandler failed to verify JoinCode"))
			transport.Forbidden(w)
			return
		}

		newNodePubKey, err := acrypto.KeyPairFromPubKeyJSON(newNodeRequest.Node.PubKey)
		if err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to KeyPairFromPubKeyJSON"))
			transport.BadRequest(w)
			return
		}

		encGlobalKey, err := newNodePubKey.Encrypt(app.KeySet.GlobalKey.JSON())
		if err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to Encrypt"))
			transport.InternalServerError(w)
			return
		}

		nodeAddedAction := actions.NewNodeAdded(newNodeRequest.Node, encGlobalKey)
		actionJSON := nodeAddedAction.JSON()

		block, err := blockchain.NewBlockWithData(app.KeySet.GlobalKey, actionJSON, nodeAddedAction.ActionType())
		if err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to NewBlockWithAction"))
			transport.InternalServerError(w)
			return
		}

		errChan := app.Chain.AddNewBlock(block)
		if err := <-errChan; err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to AddNewBlock"))
			transport.InternalServerError(w)
			return
		}

		masterPubKeyJSON := app.KeySet.KeyPair.PubKeyJSON()

		resp := requests.NewNodeResponse{
			EncGlobalKey:     encGlobalKey,
			MasterPubKeyJSON: masterPubKeyJSON,
		}

		transport.ReplyWithJSON(w, resp)
	}
}

func generateNewNodeID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return acrypto.Base64URLEncode(bytes)
}
