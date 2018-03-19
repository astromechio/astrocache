package send

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/astromechio/astrocache/model/actions"

	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/model/requests"
	mservice "github.com/astromechio/astrocache/server/master/service"
	vservice "github.com/astromechio/astrocache/server/verifier/service"
	wservice "github.com/astromechio/astrocache/server/worker/service"
	"github.com/pkg/errors"
)

// ProposeBlockToVerifiers proposes a block and decides if the verifiers will accept it
func ProposeBlockToVerifiers(block *blockchain.Block, verifiers []*model.Node, thisNode *model.Node) error {
	blockJSON, err := json.Marshal(block)
	if err != nil {
		return errors.Wrap(err, "ProposeBlockToVerifiers failed to Marshal")
	}

	req := &vservice.ProposeBlockRequest{
		Block:        blockJSON,
		ProposingNID: thisNode.NID,
	}

	// handle the single verifier case
	if verifiers == nil || len(verifiers) == 0 {
		return nil
	}

	logger.LogDebug(fmt.Sprintf("ProposeBlockToVerifiers verifying block with %d verifiers", len(verifiers)))

	resultChan := make(chan bool, len(verifiers))

	for _, v := range verifiers {
		go sendBlockProposal(v, req, resultChan)
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

	logger.LogDebug(fmt.Sprintf("ProposeBlockToVerifiers got %d matches and %d mismatches", numMatch, numMismatch))

	if numMismatch > 0 {
		return fmt.Errorf("ProposeBlockToVerifiers failed to add pending block: %d verifiers reported ID mismatch", numMismatch)
	}

	return nil
}

func sendBlockProposal(node *model.Node, req *vservice.ProposeBlockRequest, resultChan chan bool) {
	conn, err := node.Dial()
	if err != nil {
		logger.LogError(errors.Wrap(err, "sendBlockProposal failed to Dial"))
		resultChan <- false
	}

	client := vservice.NewVerifierClient(conn)

	ctx := context.Background()
	_, err = client.ProposeBlock(ctx, req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "sendBlockProposal failed to ProposeBlock"))
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
		go distributeBlock(w, req, resultChan)
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

func distributeBlock(node *model.Node, req *requests.ProposeBlockRequest, resultChan chan bool) {
	conn, err := node.Dial()
	if err != nil {
		logger.LogError(errors.Wrap(err, "distributeBlock failed to Dial"))
		resultChan <- false
	}

	blockJSON, err := json.Marshal(req.Block)
	if err != nil {
		logger.LogError(errors.Wrap(err, "distributeBlock failed to Marshal"))
		resultChan <- false
	}

	if node.Type == model.NodeTypeWorker {
		request := &wservice.AddBlockRequest{
			Block:        blockJSON,
			ProposingNID: req.ProposingNID,
		}
		client := wservice.NewWorkerClient(conn)

		ctx := context.Background()
		_, err = client.AddBlock(ctx, request)
		if err != nil {
			logger.LogError(errors.Wrap(err, "sendBlockProposal failed to ProposeBlock"))
			resultChan <- false
		}
	} else if node.Type == model.NodeTypeMaster {
		request := &mservice.AddBlockRequest{
			Block:        blockJSON,
			ProposingNID: req.ProposingNID,
		}
		client := mservice.NewMasterClient(conn)

		ctx := context.Background()
		_, err = client.AddBlock(ctx, request)
		if err != nil {
			logger.LogError(errors.Wrap(err, "sendBlockProposal failed to ProposeBlock"))
			resultChan <- false
		}
	} else {
		logger.LogError(fmt.Errorf("sendBlockProposal tried to distribute block to verifier node"))
		resultChan <- false
	}

	resultChan <- true
}
