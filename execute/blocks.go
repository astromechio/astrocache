package execute

import (
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/pkg/errors"
)

// AddPendingBlockFromRequest builds a proposed block and then attempts to add it to the chain
func AddPendingBlockFromRequest(chain *blockchain.Chain, keySet *acrypto.KeySet, req *requests.ProposeBlockRequest) (*blockchain.Block, error) {
	proposedBlock := &blockchain.Block{
		ID:         req.TempID,
		Data:       req.Data,
		ActionType: req.ActionType,
		PrevID:     req.PrevID,
	}

	err := AddPendingBlock(chain, keySet, proposedBlock)
	if err != nil {
		return nil, errors.Wrap(err, "AddPendingBlockFromRequest failed to AddPendingBlock")
	}

	return proposedBlock, nil
}

// AddPendingBlock adds a proposed block to the chain
func AddPendingBlock(chain *blockchain.Chain, keySet *acrypto.KeySet, block *blockchain.Block) error {
	err := chain.AddPendingBlock(block, keySet)
	if err != nil {
		return errors.Wrap(err, "AddPendingBlock failed to chain.AddPendingBlock")
	}

	return nil
}

// CommitPendingBlock commits a pending block to the chain
func CommitPendingBlock(chain *blockchain.Chain, keySet *acrypto.KeySet, block *blockchain.Block) error {
	return chain.CommitBlockWithTempID(block.ID, block.PrevID, nil, keySet)
}
