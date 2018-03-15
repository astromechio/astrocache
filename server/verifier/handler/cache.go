package handler

import (
	"net/http"

	"github.com/astromechio/astrocache/model/actions"
	"github.com/astromechio/astrocache/model/blockchain"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/transport"
	"github.com/pkg/errors"
)

// SetValueHandler handles value set requests
func SetValueHandler(app *config.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setValReq := &requests.SetValueRequest{}
		setValReq.FromRequest(r)

		if err := setValReq.Verify(); err != nil {
			logger.LogError(errors.Wrap(err, "SetValueHandler failed to Verify"))
			transport.BadRequest(w)
			return
		}

		action := actions.NewSetValue(setValReq.Key, setValReq.Value)
		actionJSON := action.JSON()

		block, err := blockchain.NewBlockWithData(app.KeySet.GlobalKey, actionJSON, action.ActionType())
		if err != nil {
			transport.InternalServerError(w)
			return
		}

		errChan, _ := app.Chain.ReserveBlockID(app.Self.NID)
		if err := <-errChan; err != nil {
			logger.LogError(errors.Wrap(err, "SetValueHandler failed to ReserveBlockID"))
			transport.InternalServerError(w)
			return
		}

		errChan = app.Chain.AddNewBlock(block, app.Self.NID)
		if err := <-errChan; err != nil {
			logger.LogError(errors.Wrap(err, "SetValueHandler failed to AddNewBlock"))
			transport.Conflict(w)
			return
		}

		transport.Ok(w)
	}
}
