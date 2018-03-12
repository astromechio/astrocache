package workers

import (
	"os"

	"github.com/astromechio/astrocache/model/actions"

	"github.com/astromechio/astrocache/config"

	"github.com/astromechio/astrocache/logger"
	"github.com/pkg/errors"
)

// StartActionWorker decrypts blocks and performs the actions therein
func StartActionWorker(app *config.App) {
	if app.Chain == nil {
		logger.LogError(errors.New("StartChainWorker received nil chain, terminating"))
		os.Exit(1)
	}

	if app.KeySet == nil {
		logger.LogError(errors.New("StartChainWorker received nil keySet, terminating"))
		os.Exit(1)
	}

	chain := app.Chain

	logger.LogInfo("Starting action worker")

	for true {
		block := <-chain.ActionChan

		if block == nil {
			logger.LogWarn("StartActionWorker received nil block, continuing..")
		}

		actionJSON, err := app.KeySet.GlobalKey.Decrypt(block.Data)
		if err != nil {
			logger.LogError(errors.Wrap(err, "StartActionWorker failed to Decrypt for block with ID "+block.ID))
			continue
		}

		action, err := actions.UnmarshalAction(actionJSON, block.ActionType)
		if err != nil {
			logger.LogError(errors.Wrap(err, "StartActionWorker failed to unmarshalAction for block with ID "+block.ID))
			continue
		}

		logger.LogInfo("Executing action (type " + action.ActionType() + ") from block with ID " + block.ID)

		if err := action.Execute(app); err != nil {
			logger.LogError(errors.Wrap(err, "StartActionWorker failed to Execute for block with ID "+block.ID))
			continue
		}
	}
}
