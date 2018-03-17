package service

import (
	"encoding/json"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/pkg/errors"
)

// ProposeBlockHandler handles proposed blocks
func ProposeBlockHandler(app *config.App, req *ProposeBlockRequest) (*ProposeBlockResponse, error) {
	var err error
	chain := app.Chain

	block := &blockchain.Block{}
	if err := json.Unmarshal(req.Block, block); err != nil {
		err = errors.Wrap(err, "ProposeBlockHandler failed to Unmarshal")
		logger.LogError(err)
		return nil, err
	}

	errChan := chain.VerifyProposedBlock(block, req.ProposingNID)
	if err = <-errChan; err != nil {
		err = errors.Wrap(err, "ProposeAddBlockHandler failed to AddNewBlock")
		logger.LogError(err)
		return nil, err
	}

	return &ProposeBlockResponse{}, nil
}
