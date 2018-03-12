package workers

import (
	"os"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/send"
	"github.com/pkg/errors"
)

// StartChainWorker runs on a goroutine and manages the adding of new blocks in an atomic manner
func StartChainWorker(app *config.App) {
	if app.Chain == nil {
		logger.LogError(errors.New("StartChainWorker received nil chain, terminating"))
		os.Exit(1)
	}

	if app.KeySet == nil {
		logger.LogError(errors.New("StartChainWorker received nil keySet, terminating"))
		os.Exit(1)
	}

	chain := app.Chain

	logger.LogInfo("Starting chain worker")

	for true {
		blockJob := <-chain.WorkerChan

		if blockJob.Block == nil {
			logger.LogWarn("StartChainWorker received nil block, continuing..")
			blockJob.ResultChan <- errors.New("Job had nil block")
			continue
		}

		if blockJob.Block.Signature == nil {
			if err := mineBlock(blockJob.Block, chain, app); err != nil {
				blockJob.ResultChan <- errors.Wrap(err, "StartChainWorker failed to mineBlock")
			}
		} else {
			if err := checkBlock(blockJob.Block, chain, app, blockJob.Check); err != nil {
				blockJob.ResultChan <- errors.Wrap(err, "StartChainWorker failed to checkBlock")
			}
		}

		blockJob.ResultChan <- nil
	}
}

func mineBlock(block *blockchain.Block, chain *blockchain.Chain, app *config.App) error {
	logger.LogInfo("mineBlock mining block")

	block.PrepareForCommit(app.KeySet.KeyPair, chain.LastBlock())

	if err := chain.SetProposedBlock(block); err != nil {
		return errors.Wrap(err, "mineBlock failed to SetProposedBlock")
	}

	if err := send.ProposeBlockToVerifiers(block, app.NodeList.Verifiers, app.Self); err != nil {
		return errors.Wrap(err, "mineBlock failed to ProposeBlockToVerifiers")
	}

	if err := chain.CommitProposedBlock(app.KeySet); err != nil {
		return errors.Wrap(err, "mineBlock failed to CommitProposedBlock")
	}

	return nil
}

func checkBlock(block *blockchain.Block, chain *blockchain.Chain, app *config.App, check bool) error {
	logger.LogInfo("checkBlock checking block")

	if err := chain.SetProposedBlock(block); err != nil {
		return errors.Wrap(err, "checkBlock failed to SetProposedBlock")
	}

	if check {
		if err := send.CheckBlockWithVerifiers(block, app.NodeList.Verifiers, app.Self); err != nil {
			return errors.Wrap(err, "checkBlock failed to CheckBlockWithVerifiers")
		}
	}

	if err := chain.CommitProposedBlock(app.KeySet); err != nil {
		return errors.Wrap(err, "checkBlock failed to CommitProposedBlock")
	}

	return nil
}
