package handler

import (
	"crypto/rand"
	"net/http"

	"github.com/astromechio/astrocache/model/blockchain"

	"github.com/astromechio/astrocache/logger"
	"github.com/pkg/errors"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/actions"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/server"
	"github.com/astromechio/astrocache/transport"
)

// AddVerifierNodeHandler handles POST /v1/master/nodes/verifier
func AddVerifierNodeHandler(config *server.Config) http.HandlerFunc {
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

		joinCode := config.ValueForKey(server.ConfigJoinCodeKey)
		if newNodeRequest.JoinCode != joinCode {
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

		block, err := blockchain.NewBlockWithAction(config.KeySet.GlobalKey, nodeAddedAction)
		if err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to NewBlockWithAction"))
			transport.InternalServerError(w)
			return
		}

		errChan := config.Chain.AddNewBlock(block)
		if err := <-errChan; err != nil {
			logger.LogError(errors.Wrap(err, "AddVerifierNodeHandler failed to AddNewBlock"))
			transport.InternalServerError(w)
			return
		}

		config.NodeList.AddVerifier(newNode)

		masterPubKeyJSON := config.KeySet.KeyPair.PubKeyJSON()

		resp := requests.NewNodeResponse{
			EncGlobalKey:     encGlobalKey,
			Node:             newNode,
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
