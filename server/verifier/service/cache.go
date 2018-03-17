package service

import (
	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/actions"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/pkg/errors"
)

// SetValueHandler handles set value requests
func SetValueHandler(app *config.App, req *SetValueRequest) (*SetValueResponse, error) {
	var err error

	action := actions.NewSetValue(req.Key, req.Value)
	actionJSON := action.JSON()

	block, err := blockchain.NewBlockWithData(app.KeySet.GlobalKey, actionJSON, action.ActionType())
	if err != nil {
		err = errors.Wrap(err, "SetValueHandler failed to NewBlockWithData")
		logger.LogError(err)
		return nil, err
	}

	errChan, _ := app.Chain.ReserveBlockID(app.Self.NID)
	if err := <-errChan; err != nil {
		err = errors.Wrap(err, "SetValueHandler failed to ReserveBlockID")
		logger.LogError(err)
		return nil, err
	}

	errChan = app.Chain.AddNewBlock(block, app.Self.NID)
	if err := <-errChan; err != nil {
		err = errors.Wrap(err, "SetValueHandler failed to AddNewBlock")
		logger.LogError(err)
		return nil, err
	}

	return &SetValueResponse{}, nil
}
