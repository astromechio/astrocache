package execute

import (
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/pkg/errors"
)

// AddPendingBlock adds a proposed block to the chain
func AddPendingBlock(chain *blockchain.Chain, keySet *acrypto.KeySet, req *requests.ProposeBlockRequest) (*requests.ProposeBlockResponse, error) {
	proposedBlock := &blockchain.Block{
		ID:         req.TempID,
		Data:       req.Data,
		ActionType: req.ActionType,
	}

	prevBlock := chain.AddPendingBlock(proposedBlock, keySet)

	prevHash, err := prevBlock.Hash()
	if err != nil {
		return nil, errors.Wrap(err, "AddPendingBlock failed to prevBlock.Hash")
	}

	res := &requests.ProposeBlockResponse{
		PrevID:   prevBlock.ID,
		PrevHash: prevHash,
	}

	return res, nil
}
