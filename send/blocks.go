package send

import (
	"bytes"
	"fmt"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/execute"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/actions"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/transport"
	"github.com/pkg/errors"
)

// ProposeBlockWithAction proposes a block and decides if the verifiers will accept it
func ProposeBlockWithAction(verifiers []*model.Node, action actions.Action, keySet *acrypto.KeySet, chain *blockchain.Chain) (*blockchain.Block, error) {
	block, err := blockchain.NewBlockWithAction(keySet.GlobalKey, action)
	if err != nil {
		return nil, errors.Wrap(err, "ProposeBlockWithAction failed to NewBlockWithAction")
	}

	prevBlock := chain.LastBlock()
	block.PrevID = prevBlock.ID

	req := &requests.ProposeBlockRequest{
		TempID:     block.ID,
		Data:       block.Data,
		ActionType: action.ActionType(),
		PrevID:     block.PrevID,
	}

	resultChan := make(chan bool)

	prevHash, err := prevBlock.Hash()
	if err != nil {
		return nil, errors.Wrap(err, "ProposeBlockWithAction failed to prevBlock.Hash()")
	}

	for _, v := range verifiers {
		reqURL := transport.URLFromAddressAndPath(v.Address, req.Path())

		go sendBlockProposal(reqURL, req, prevHash, resultChan)
	}

	numMatch := 0
	numMismatch := 0

	for true {
		match := <-resultChan

		if match {
			numMatch++
		} else {
			numMismatch++
		}

		if numMatch+numMismatch == len(verifiers) {
			break
		}
	}

	if numMismatch > 0 {
		return nil, fmt.Errorf("ProposeBlockWithAction failed to add pending block: %d verifiers reported ID mismatch", numMismatch)
	}

	res, err := execute.AddPendingBlock(chain, keySet, req)
	if err != nil {
		if res != nil {
			logger.LogError(errors.New("ProposeBlockWithAction tried to AddPendingBlock, but an ID mismatch occurred"))
		}

		return nil, errors.Wrap(err, "ProposeBlockWithAction failed to AddPendingBlock")
	}

	return block, nil
}

func sendBlockProposal(url string, req *requests.ProposeBlockRequest, prevHash []byte, resultChan chan bool) {
	resp := &requests.ProposeBlockResponse{}

	if err := transport.Post(url, req, resp); err != nil {
		logger.LogError(errors.Wrap(err, "sendBlockProposal failed to Post"))
		resultChan <- false
	}

	if !bytes.Equal(prevHash, resp.PrevHash) {
		logger.LogError(errors.New("sendBLockProposal found prevHash mismatch"))
		resultChan <- false
		return
	}

	resultChan <- resp.IDMismatch
}
