package workers

import (
	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/send"
	"github.com/pkg/errors"
)

// DistributeWorker starts the distribute worker
func DistributeWorker(app *config.App) {
	chain := app.Chain

	logger.LogInfo("starting distribute worker")

	for true {
		block := <-chain.DistributeChan

		// update this every time in case we got a new worker
		workers := app.NodeList.WorkersForVerifierWithNID(app.Self.NID)

		if err := send.DistributeBlockToWorkers(block, workers, app.Self); err != nil {
			logger.LogError(errors.Wrap(err, "DistributeWorker failed to DistributeBlockToWorkers"))
		}
	}
}
