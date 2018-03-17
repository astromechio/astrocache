package service

import (
	"encoding/json"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/pkg/errors"
)

// AddBlockHandler handles adding blocks
func AddBlockHandler(app *config.App, req *AddBlockRequest) (*AddBlockResponse, error) {
	var err error
	chain := app.Chain

	block := &blockchain.Block{}
	if err := json.Unmarshal(req.Block, block); err != nil {
		err = errors.Wrap(err, "AddBlockHandler failed to Unmarshal")
		return nil, err
	}

	errChan := chain.VerifyProposedBlock(block, "")
	if err = <-errChan; err != nil {
		err = errors.Wrap(err, "AddBlockHandler failed to AddNewBlock")
		logger.LogError(err)
		return nil, err
	}

	return &AddBlockResponse{}, nil
}
