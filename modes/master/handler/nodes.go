package handler

import (
	"crypto/rand"
	"net/http"

	"github.com/astromechio/astrocache/execute"

	"github.com/astromechio/astrocache/logger"
	"github.com/pkg/errors"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/actions"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/modes/master/config"
	"github.com/astromechio/astrocache/transport"
)

// AddVerifierNodeHandler handles POST /v1/master/nodes/verifier
func AddVerifierNodeHandler(config *config.Config) http.HandlerFunc {
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

		if newNodeRequest.JoinCode != config.JoinCode {
			logger.LogError(errors.New("AddVerifierNodeHandler failed to verify JoinCode"))
			transport.Forbidden(w)
			return
		}

		newNode := &model.Node{
			NID:     generateNewNodeID(),
			Address: newNodeRequest.Address,
			Type:    model.NodeTypeVerifier,
			PubKey:  newNodeRequest.PubKey,
		}

		newNodePubKey, err := acrypto.KeyPairFromPubKeyJSON(newNodeRequest.PubKey)
		if err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to KeyPairFromPubKeyJSON"))
			transport.BadRequest(w)
			return
		}

		encGlobalKey, err := newNodePubKey.Encrypt(config.KeySet.GlobalKey.JSON())
		if err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to Encrypt"))
			transport.InternalServerError(w)
			return
		}

		nodeAddedAction := actions.NewNodeAdded(newNode, encGlobalKey)

		// in the bootstrap case, verNode will be nil, which is fine
		verNode := config.Nodes.RandomVerifier()

		resp, err := execute.AddNodeToChain(config.Chain, config.KeySet, verNode, nodeAddedAction)
		if err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to AddNodeToChain"))
			transport.InternalServerError(w)
			return
		}

		transport.ReplyWithJSON(w, resp)
	}
}

func generateNewNodeID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return acrypto.Base64URLEncode(bytes)
}
