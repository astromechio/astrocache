package workers

import (
	"fmt"
	"os"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
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
		var reserveJob *blockchain.ReserveIDJob
		var blockJob *blockchain.NewBlockJob

		// if we're a worker, we don't care about reservations
		if app.Self.Type != model.NodeTypeWorker {
			select {
			case reserveJob = <-chain.ReservedChan:
				logger.LogInfo("PoposeWorker got reserved block job")

				blockJob = <-chain.ProposeChan
				logger.LogInfo("ProposeWorker got proposed block job")

				if err := proposeBlock(blockJob, app); err != nil {
					blockJob.ResultChan <- errors.Wrap(err, "ProposeWorker failed to proposeBlock")
					chain.Proposed = nil

					continue
				}

				if blockJob.Block.ID != reserveJob.BlockID {
					blockJob.ResultChan <- fmt.Errorf("ProposeWorker failed after proposeBlock, block.ID did not match reserved.BlockID")
					chain.Proposed = nil

					continue
				}

			case blockJob = <-chain.VerifyChan:
				logger.LogInfo("ProposeWorker got verify block job")

				if reserveJob != nil && blockJob.Block.ID != reserveJob.BlockID {
					logger.LogError(fmt.Errorf("ProposeWorker failed before verifyBlock, block.ID did not match reserved.BlockID"))
					blockJob.ResultChan <- fmt.Errorf("ProposeWorker failed before verifyBlock, block.ID did not match reserved.BlockID")
					chain.Proposed = nil

					continue
				}

				chain.Proposed = blockJob.Block
			}
		} else {
			blockJob = <-chain.VerifyChan
			chain.Proposed = blockJob.Block
		}

		// Send the block to be committed
		logger.LogInfo("ProposeWorker sending blockJob to be committed")

		select {
		case <-chain.CommittedChan:
			logger.LogInfo("")
		default:
			// we want to clear out CommittedChan in case the ReserveWorker is busy
		}

		chain.CommitChan <- blockJob
		logger.LogInfo("ProposeWorker sent blockJob to be committed")

		if app.Self.Type == model.NodeTypeWorker {
			// workers will wait in the committed loop
			<-chain.CommittedChan
		}
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

	if err := send.ProposeBlockToVerifiers(job.Block, app.NodeList.Verifiers, app.Self); err != nil {
		return errors.Wrap(err, "proposeBlock failed to ProposeBlockToVerifiers")
	}

	chain.Proposed = job.Block

	return nil
}
