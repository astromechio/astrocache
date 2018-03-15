package workers

import (
	"fmt"
	"os"
	"time"

	"github.com/astromechio/astrocache/model/blockchain"

	"github.com/astromechio/astrocache/config"
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/send"
	"github.com/pkg/errors"
)

// ReserveWorker runs on a goroutine and manages the adding of new blocks in an atomic manner
func ReserveWorker(app *config.App) {
	if app.Chain == nil {
		logger.LogError(errors.New("ReserveWorker received nil chain, terminating"))
		os.Exit(1)
	}

	chain := app.Chain

	logger.LogInfo("starting reserve worker")

	for true {
		var reserveJob *blockchain.ReserveIDJob

		select {
		case reserveJob = <-chain.ReserveChan:
			logger.LogInfo("ReserveWorker got reserve job")
		case <-chain.CommittedChan:
			logger.LogInfo("ReserveWorker saw committed message, continuing...")
			continue
		}

		if app.Self.Type == model.NodeTypeVerifier {
			reservedID, err := send.RequestReservedID(app.NodeList.Master, app.Self.NID)
			if err != nil {
				logger.LogError(errors.Wrap(err, "ReserveWorker failed to RequesReservedID"))
				reserveJob.ResultChan <- errors.Wrap(err, "ReserveWorker failed to RequesReservedID")
				continue
			}

			reserveJob.BlockID = reservedID.BlockID

			logger.LogInfo(fmt.Sprintf("ReserveWorker got reserved block ID %s, sending reserved message", reserveJob.BlockID))
			chain.ReservedChan <- reserveJob

			logger.LogInfo("ReserveWorker sending job result")
			reserveJob.ResultChan <- nil

			logger.LogInfo("ReserveWorker finished job, waiting for committed")
			<-chain.CommittedChan

		} else {
			last := chain.LastBlock()

			hash, err := last.Hash()
			if err != nil {
				logger.LogError(errors.Wrap(err, "ReserveWorker failed to last.Hash()"))
				reserveJob.ResultChan <- errors.Wrap(err, "ReserveWorker failed to last.Hash()")
				continue
			}

			newID := acrypto.Base64URLEncode(hash)

			reserveJob.BlockID = newID

			logger.LogInfo(fmt.Sprintf("ReserveWorker reserved block with ID %s", newID))

			if reserveJob.ProposingNID == app.Self.NID {
				// if we are mining our own block
				logger.LogInfo("ReserveWorker finished job with propNID of self, sending reserved message")
				chain.ReservedChan <- reserveJob

				logger.LogInfo("ReserveWorker sending job result")
				reserveJob.ResultChan <- nil
			} else {
				logger.LogInfo("ReserveWorker sending job result")
				reserveJob.ResultChan <- nil

				logger.LogInfo("ReserveWorker finished job, waiting for committed or timeout")
				select {
				case <-chain.CommittedChan:
				case <-time.After(time.Second * 2):
					logger.LogInfo("ReserveWorker hit timeout waiting for committed, assuming failed")
				}
			}
		}
	}
}
