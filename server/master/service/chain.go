package service

import (
	"encoding/json"
	fmt "fmt"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/pkg/errors"
)

// GetChainHandler handles the GetChain rpc for a master node
func GetChainHandler(app *config.App, req *GetChainRequest) (*GetChainResponse, error) {
	var chainJSON []byte
	var err error

	if req.After == "" {
		chainJSON, err = json.Marshal(app.Chain.Blocks)
	} else {
		blocks := app.Chain.BlocksAfterID(req.After)
		if blocks == nil {
			return nil, fmt.Errorf("GetChainHandler failed to BlocksAfterID: block with ID %s does not exist", req.After)
		}

		chainJSON, err = json.Marshal(blocks)
		if err != nil {
			return nil, errors.Wrap(err, "GetChainHandler failed to Marshal")
		}
	}

	res := &GetChainResponse{
		Blocks: chainJSON,
	}

	return res, nil
}

// ReserveIDHandler handles the block reservation rpc for a master node
func ReserveIDHandler(app *config.App, req *ReserveIDRequest) (*ReserveIDResponse, error) {
	var err error

	errChan, reserveJob := app.Chain.ReserveBlockID(req.ProposingNID)
	if err = <-errChan; err != nil {
		err = errors.New("ReserveIDHandler failed to ReserveBlockID, got empty blockID")
		logger.LogError(err)
		return nil, err
	}

	resp := &ReserveIDResponse{
		BlockID: reserveJob.BlockID,
	}

	return resp, nil
}
