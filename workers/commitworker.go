package workers

import (
	"fmt"
	"os"
	"time"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/send"
	"github.com/pkg/errors"
)

// CommitWorker runs on a goroutine and manages the adding of new blocks in an atomic manner
func CommitWorker(app *config.App) {
	if app.Chain == nil {
		logger.LogError(errors.New("CommitWorker received nil chain, terminating"))
		os.Exit(1)
	}

	if app.KeySet == nil {
		logger.LogError(errors.New("CommitWorker received nil keySet, terminating"))
		os.Exit(1)
	}

	chain := app.Chain

	logger.LogInfo("starting commit worker")

	for true {
		blockJob := <-chain.CommitChan

		if blockJob.Block == nil {
			logger.LogWarn("CommitWorker received nil block, continuing..")
			continue
		}

		if err := commitBlock(blockJob, app); err != nil {
			if blockJob.Check == false {
				blockJob.ResultChan <- errors.Wrap(err, "CommitWorker failed to checkBlock")
			} else {
				logger.LogError(errors.Wrap(err, "CommitWorker failed to checkBlock and blockJob.ResultChan is closed"))
			}

			chain.Proposed = nil

			chain.CommittedChan <- nil
			loadMissingBlocks(app)

			continue
		}

		// if check is false, then we are the proposer or there was no check to begin with
		if blockJob.Check == false {
			blockJob.ResultChan <- nil
		}

		chain.ActionChan <- blockJob.Block // send the block to be executed

		chain.CommittedChan <- blockJob.Block // notify other goroutines that something was committed
	}
}

func commitBlock(job *blockchain.NewBlockJob, app *config.App) error {
	chain := app.Chain

	logger.LogInfo(fmt.Sprintf("commitBlock committing block with ID %q", job.Block.ID))

	if !job.Block.IsSameAsBlock(chain.Proposed) {
		return fmt.Errorf("commitBlock tried to commit a non-proposed block")
	}

	last := chain.LastBlock()
	if last != nil && job.Block.IsSameAsBlock(last) {
		logger.LogWarn("Tried committing duplicate block, skipping...")
		return nil
	}

	prevBlock := chain.LastBlock()

	// Verify handles the genesis case
	if err := job.Block.Verify(app.KeySet, prevBlock); err != nil {
		return errors.Wrap(err, "commitBlock failed to block.Verify")
	}

	if job.Check {
		if err := send.CheckBlockWithVerifiers(job.Block, app.NodeList.Verifiers, job.ProposingNID); err != nil {
			return errors.Wrap(err, "commitBlock failed to CheckBlockWithWorkers")
		}
	}

	logger.LogInfo(fmt.Sprintf("*** Committing bock with ID %q ***", job.Block.ID))

	chain.Blocks = append(chain.Blocks, chain.Proposed)

	chain.Proposed = nil

	return nil
}

func loadMissingBlocks(app *config.App) {
	lastBlock := app.Chain.LastBlock()

	logger.LogInfo(fmt.Sprintf("loadMissingBlocks attempting to load missing blocks after %q", lastBlock.ID))

	var missing []*blockchain.Block

	for true {
		var err error
		missing, err = send.GetBlocksAfter(app.NodeList.Master, lastBlock.ID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "loadMissingBlocks failed to GetBlocksAfter"))
			return
		}

		if missing == nil {
			logger.LogInfo(fmt.Sprintf("loadMissingBlocks received no blocks after %q", lastBlock.ID))
		} else if len(missing) > 0 {
			break
		}

		<-time.After(time.Second * 1)
	}

	for i := range missing {
		logger.LogInfo(fmt.Sprintf("loadMissingBlocks loading missing block with ID %q", missing[i].ID))

		errChan := app.Chain.VerifyBlockUnchecked(missing[i])
		if err := <-errChan; err != nil {
			logger.LogError(errors.Wrap(err, "loadMissingBlocks failed to AddNewBlockUnchecked for block with ID "+missing[i].ID))
			return
		}
	}

	logger.LogInfo(fmt.Sprintf("loadMissingBlocks loaded %d missing blocks", len(missing)))
}
