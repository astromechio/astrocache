package service

import (
	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/model"
	context "golang.org/x/net/context"
)

// MasterService exposes a master service
type MasterService struct {
	App *config.App
}

// RequestReservedID handles ID reservations
func (ms *MasterService) RequestReservedID(ctx context.Context, in *ReserveIDRequest) (*ReserveIDResponse, error) {
	return ReserveIDHandler(ms.App, in)
}

// AddNode handles adding nodes
func (ms *MasterService) AddNode(ctx context.Context, in *NewNodeRequest) (*NewNodeResponse, error) {
	return AddNodeHandler(ms.App, in)
}

// GetChain handles getting chain
func (ms *MasterService) GetChain(ctx context.Context, in *GetChainRequest) (*GetChainResponse, error) {
	return GetChainHandler(ms.App, in)
}

// AddBlock handles adding a block
func (ms *MasterService) AddBlock(ctx context.Context, in *AddBlockRequest) (*AddBlockResponse, error) {
	return AddBlockHandler(ms.App, in)
}

// ProtoNodeFromNode converts a model.Node to a proto.Node
func ProtoNodeFromNode(n *model.Node) *Node {
	return &Node{
		NID:       n.NID,
		Address:   n.Address,
		Type:      n.Type,
		PubKey:    n.PubKey,
		ParentNID: n.ParentNID,
	}
}

// NodeFromProtoNode converts a model.Node to a proto.Node
func NodeFromProtoNode(n *Node) *model.Node {
	return &model.Node{
		NID:       n.NID,
		Address:   n.Address,
		Type:      n.Type,
		PubKey:    n.PubKey,
		ParentNID: n.ParentNID,
	}
}
