package send

import (
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/actions"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/pkg/errors"
)

// AddNodeToChain handles adding a new node to the network
func AddNodeToChain(chain *blockchain.Chain, keySet *acrypto.KeySet, verifier *model.Node, action *actions.NodeAdded) error {
	actionJSON := action.JSON()

	block, err := blockchain.NewBlockWithData(keySet.GlobalKey, actionJSON, action.ActionType())
	if err != nil {
		return errors.Wrap(err, "AddNodeToChain failed to NewBlockWithAction")
	}

	prevBlock := chain.LastBlock()
	block.PrevID = prevBlock.ID

	// handle the bootstrap case
	if verifier == nil {
		errChan := chain.AddNewBlock(block)
		if err := <-errChan; err != nil {
			return errors.Wrap(err, "AddNodeToChain failed to AddNewBlock")
		}
	} else {
		logger.LogWarn("Verifier node block mining not yet implemented")
	}

	return nil
}
