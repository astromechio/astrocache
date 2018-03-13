package workers

import (
	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/send"
	"github.com/pkg/errors"
)

// StartDistributeWorker starts the distribute worker
func StartDistributeWorker(app *config.App) {
	chain := app.Chain

	logger.LogInfo("Starting distribute worker")

	for true {
		block := <-chain.DistributeChan

		// update this every time in case we got a new worker
		workers := app.NodeList.WorkersForVerifierWithNID(app.Self.NID)

		if err := send.DistributeBlockToWorkers(block, workers, app.Self); err != nil {
			logger.LogError(errors.Wrap(err, "StartDistributeWorker failed to DistributeBlockToWorkers"))
		}
	}
}
