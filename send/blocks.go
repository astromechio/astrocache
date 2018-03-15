package send

import (
	"fmt"

	"github.com/astromechio/astrocache/model/actions"

	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/transport"
	"github.com/pkg/errors"
)

// ProposeBlockToVerifiers proposes a block and decides if the verifiers will accept it
func ProposeBlockToVerifiers(block *blockchain.Block, verifiers []*model.Node, thisNode *model.Node) error {
	req := &requests.ProposeBlockRequest{
		Block:        block,
		ProposingNID: thisNode.NID,
	}

	// handle the single verifier case
	if verifiers == nil || len(verifiers) == 0 {
		return nil
	}

	logger.LogInfo(fmt.Sprintf("ProposeBlockToVerifiers verifying block with %d verifiers", len(verifiers)))

	resultChan := make(chan bool, len(verifiers))

	for _, v := range verifiers {
		reqURL := transport.URLFromAddressAndPath(v.Address, req.Path())

		go sendBlockProposal(reqURL, req, resultChan)
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

	logger.LogInfo(fmt.Sprintf("ProposeBlockToVerifiers got %d matches and %d mismatches", numMatch, numMismatch))

	if numMismatch > 0 {
		return fmt.Errorf("ProposeBlockToVerifiers failed to add pending block: %d verifiers reported ID mismatch", numMismatch)
	}

	return nil
}

func sendBlockProposal(url string, req *requests.ProposeBlockRequest, resultChan chan bool) {
	if err := transport.Post(url, req, nil); err != nil {
		logger.LogError(errors.Wrap(err, "sendBlockProposal failed to Post"))
		resultChan <- false
	}

	resultChan <- true
}

// CheckBlockWithVerifiers proposes a block and decides if the verifiers will accept it
func CheckBlockWithVerifiers(block *blockchain.Block, verifiers []*model.Node, propNID string) error {
	req := &requests.CheckBlockRequest{
		Block: block,
	}

	actualVerifiers := []*model.Node{}
	for i, v := range verifiers {
		if v.NID != propNID {
			actualVerifiers = append(actualVerifiers, verifiers[i])
		}
	}

	// handle the single verifier case
	if len(actualVerifiers) == 0 {
		return nil
	}

	logger.LogInfo(fmt.Sprintf("CheckBlockWithVerifiers checking block with %d verifiers and propNID %q", len(actualVerifiers), propNID))

	resultChan := make(chan bool)

	for _, v := range actualVerifiers {
		reqURL := transport.URLFromAddressAndPath(v.Address, req.Path())

		go sendBlockCheck(reqURL, req, resultChan)
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

		if numMatch+numMismatch == len(actualVerifiers) {
			break
		}
	}

	logger.LogInfo(fmt.Sprintf("CheckBlockWithVerifiers got %d matches and %d mismatches", numMatch, numMismatch))

	if numMismatch > 0 {
		return fmt.Errorf("CheckBlockWithVerifiers failed to check pending block: %d verifiers reported ID mismatch", numMismatch)
	}

	return nil
}

func sendBlockCheck(url string, req *requests.CheckBlockRequest, resultChan chan bool) {
	if err := transport.Post(url, req, nil); err != nil {
		logger.LogError(errors.Wrap(err, "sendBlockCheck failed to Post"))
		resultChan <- false
	}

	resultChan <- true
}

// DistributeBlockToWorkers sends a block to all workers
func DistributeBlockToWorkers(block *blockchain.Block, workers []*model.Node, thisNode *model.Node) error {
	req := &requests.ProposeBlockRequest{
		Block:        block,
		ProposingNID: thisNode.NID,
	}

	realWorkers := []*model.Node{}

	// we don't want to send node added blocks back to the master, so we
	// have to be hacky and remove the maste node
	if block.ActionType == actions.ActionTypeNodeAdded {
		for i, w := range workers {
			if w.Type != model.NodeTypeMaster {
				realWorkers = append(realWorkers, workers[i])
			}
		}
	} else {
		realWorkers = workers
	}

	if len(realWorkers) == 0 {
		return nil
	}

	logger.LogInfo(fmt.Sprintf("DistributeBlockToWorkers distributing block with ID %q to %d workers", block.ID, len(realWorkers)))

	resultChan := make(chan bool)

	numFailed := 0
	numSucceeded := 0

	for _, w := range realWorkers {
		reqURL := transport.URLFromAddressAndPath(w.Address, "v1/worker/block")

		go distributeBlock(reqURL, req, resultChan)
	}

	for true {
		succeeded := <-resultChan

		if succeeded {
			numSucceeded++
		} else {
			numFailed++
		}

		if numSucceeded+numFailed == len(realWorkers) {
			break
		}
	}

	if numFailed > 0 {
		return fmt.Errorf("DistributeBlockToWorkers failed to distribute block with ID %q to %d workers", block.ID, numFailed)
	}

	return nil
}

func distributeBlock(url string, req *requests.ProposeBlockRequest, resultChan chan bool) {
	if err := transport.Post(url, req, nil); err != nil {
		logger.LogError(errors.Wrap(err, "distributeBlock failed to Post"))
		resultChan <- false
	}

	resultChan <- true
}
