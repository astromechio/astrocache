package workers

import (
	"fmt"
	"os"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/send"
	"github.com/pkg/errors"
)

// ProposeWorker runs on a goroutine and manages the adding of new blocks in an atomic manner
func ProposeWorker(app *config.App) {
	if app.Chain == nil {
		logger.LogError(errors.New("ProposeWorker received nil chain, terminating"))
		os.Exit(1)
	}

	if app.KeySet == nil {
		logger.LogError(errors.New("ProposeWorker received nil keySet, terminating"))
		os.Exit(1)
	}

	chain := app.Chain

	logger.LogInfo("starting propose worker")

	for true {
		logger.LogInfo(fmt.Sprintf("ProposeWorker shows %d jobs in ProposeChan", len(chain.ProposeChan)))
		logger.LogInfo(fmt.Sprintf("ProposeWorker shows %d jobs in VerifyChan", len(chain.VerifyChan)))
		var blockJob *blockchain.NewBlockJob

		select {
		case blockJob = <-chain.VerifyChan:
			if err := verifyBlock(blockJob.Block, app); err != nil {
				blockJob.ResultChan <- errors.Wrap(err, "ProposeWorker failed to verifyBlock")
				chain.Proposed = nil

				continue
			}

			// if check is true, we have a handler waiting for this response so we return the result now
			if blockJob.Check == true {
				blockJob.ResultChan <- nil
			}

		case blockJob = <-chain.ProposeChan:
			if err := proposeBlock(blockJob, app); err != nil {
				blockJob.ResultChan <- errors.Wrap(err, "ProposeWorker failed to proposeBlock")
				chain.Proposed = nil

				continue
			}

			blockJob.Check = false // since we've already verified the block with the verifiers, it doesn't need to be checked again

		}

		// Send the block to be committed
		chain.CommitChan <- blockJob

		// Notify other goroutines that something was proposed
		//chain.ProposedChan <- blockJob.Block

		<-chain.CommittedChan // wait until the next block is committed
	}
}

func proposeBlock(job *blockchain.NewBlockJob, app *config.App) error {
	chain := app.Chain

	prevBlock := chain.LastBlock()
	job.Block.PrepareForCommit(app.KeySet.KeyPair, prevBlock)

	logger.LogInfo(fmt.Sprintf("proposeBlock proposing block with ID %q", job.Block.ID))

	// TODO: determine if this is needed
	if err := job.Block.Verify(app.KeySet, prevBlock); err != nil {
		return errors.Wrap(err, "proposeBlock failed to Verify")
	}

	if job.Check {
		// ask the other verifier nodes if this block is good
		if err := send.ProposeBlockToVerifiers(job.Block, app.NodeList.Verifiers, app.Self); err != nil {
			return errors.Wrap(err, "proposeBlock failed to ProposeBlockToVerifiers")
		}
	}

	chain.Proposed = job.Block

	return nil
}

func verifyBlock(block *blockchain.Block, app *config.App) error {
	chain := app.Chain

	logger.LogInfo(fmt.Sprintf("verifyBlock checking block with ID %q", block.ID))

	prevBlock := chain.LastBlock()

	if err := block.Verify(app.KeySet, prevBlock); err != nil {
		return errors.Wrap(err, "verifyBlock failed to Verify")
	}

	chain.Proposed = block

	return nil
}
