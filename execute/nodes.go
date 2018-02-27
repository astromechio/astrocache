package execute

import (
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/actions"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/pkg/errors"
)

// AddNodeToChain handles adding a new node to the network
func AddNodeToChain(chain *blockchain.Chain, keySet *acrypto.KeySet, verifier *model.Node, action *actions.NodeAdded) (*requests.NewNodeResponse, error) {
	block, err := blockchain.NewBlockWithAction(keySet.GlobalKey, action)
	if err != nil {
		return nil, errors.Wrap(err, "AddNodeToChain failed to NewBlockWithAction")
	}

	// handle the bootstrap case
	if verifier == nil {
		prevBlock := chain.AddPendingBlock(block, keySet)

		if err := chain.CommitBlockWithTempID(block.ID, prevBlock.ID, nil, keySet); err != nil {
			return nil, errors.Wrap(err, "AddNodeToChain failed to CommitBlockWithTempID")
		}
	} else {
		logger.LogWarn("Verifier node block mining not yet implemented")
	}

	resp := &requests.NewNodeResponse{
		Node:         action.Node,
		EncGlobalKey: action.EncGlobalKey,
	}

	return resp, nil
}
