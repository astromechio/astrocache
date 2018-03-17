package send

import (
	"context"
	"encoding/json"

	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/model/requests"
	mservice "github.com/astromechio/astrocache/server/master/service"
	"github.com/pkg/errors"
)

// GetEntireChain requests the entire chain from the master node
func GetEntireChain(master *model.Node) ([]*blockchain.Block, error) {
	return GetBlocksAfter(master, "")
}

// GetBlocksAfter requests the entire chain from the master node
func GetBlocksAfter(master *model.Node, afterID string) ([]*blockchain.Block, error) {
	req := &mservice.GetChainRequest{
		After: afterID,
	}

	conn, err := master.Dial()
	if err != nil {
		return nil, errors.Wrap(err, "JoinNetwork failed to Dial")
	}

	client := mservice.NewMasterClient(conn)

	ctx := context.Background()
	resp, err := client.GetChain(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "GetBlocksAfter failed to client.GetChain")
	}

	blocks := []*blockchain.Block{}
	if err := json.Unmarshal(resp.Blocks, &blocks); err != nil {
		return nil, errors.Wrap(err, "GetBlocksAfter failed to Unmarshal")
	}

	return blocks, nil
}

// RequestReservedID reserves a block ID with the master node
func RequestReservedID(master *model.Node, propNID string) (*requests.ReserveIDResponse, error) {
	logger.LogInfo("RequestReservedID requesting block ID from master node")

	req := &mservice.ReserveIDRequest{
		ProposingNID: propNID,
	}

	conn, err := master.Dial()
	if err != nil {
		return nil, errors.Wrap(err, "JoinNetwork failed to Dial")
	}

	client := mservice.NewMasterClient(conn)

	ctx := context.Background()
	resp, err := client.RequestReservedID(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "RequestReservedID failed to client.RequestReservedID")
	}

	response := &requests.ReserveIDResponse{
		BlockID: resp.BlockID,
	}

	return response, nil
}
