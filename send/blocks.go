package send

import (
	"fmt"

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
		Block:    block,
		MinerNID: thisNode.NID,
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
func CheckBlockWithVerifiers(block *blockchain.Block, verifiers []*model.Node, thisNode *model.Node) error {
	req := &requests.CheckBlockRequest{
		Block: block,
	}

	// handle the single verifier case
	if len(verifiers) == 0 {
		return nil
	}

	resultChan := make(chan bool)

	for _, v := range verifiers {
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

		if numMatch+numMismatch == len(verifiers) {
			break
		}
	}

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
