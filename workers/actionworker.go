package workers

import (
	"os"

	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/actions"

	"github.com/astromechio/astrocache/config"

	"github.com/astromechio/astrocache/logger"
	"github.com/pkg/errors"
)

// ActionWorker decrypts blocks and performs the actions therein
func ActionWorker(app *config.App) {
	if app.Chain == nil {
		logger.LogError(errors.New("StartChainWorker received nil chain, terminating"))
		os.Exit(1)
	}

	if app.KeySet == nil {
		logger.LogError(errors.New("StartChainWorker received nil keySet, terminating"))
		os.Exit(1)
	}

	chain := app.Chain

	logger.LogInfo("starting action worker")

	for true {
		block := <-chain.ActionChan

		if block == nil {
			logger.LogWarn("ActionWorker received nil block, continuing..")
		}

		actionJSON, err := app.KeySet.GlobalKey.Decrypt(block.Data)
		if err != nil {
			logger.LogError(errors.Wrap(err, "ActionWorker failed to Decrypt for block with ID "+block.ID))
			continue
		}

		action, err := actions.UnmarshalAction(actionJSON, block.ActionType)
		if err != nil {
			logger.LogError(errors.Wrap(err, "ActionWorker failed to unmarshalAction for block with ID "+block.ID))
			continue
		}

		logger.LogDebug("executing action (type " + action.ActionType() + ") from block with ID " + block.ID)

		if err := action.Execute(app); err != nil {
			logger.LogError(errors.Wrap(err, "ActionWorker failed to Execute for block with ID "+block.ID))
			continue
		} else {
			if app.Self.Type == model.NodeTypeVerifier {
				// distribute the block to all worker nodes we care about
				chain.DistributeChan <- block
			}
		}
	}
}
