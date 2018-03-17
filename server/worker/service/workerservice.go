package service

import (
	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/model"
	context "golang.org/x/net/context"
)

// WorkerService exposes a master service
type WorkerService struct {
	App *config.App
}

// AddBlock handles adding a block
func (ms *WorkerService) AddBlock(ctx context.Context, in *AddBlockRequest) (*AddBlockResponse, error) {
	return AddBlockHandler(ms.App, in)
}

// SetValue handles adding a block
func (ms *WorkerService) SetValue(ctx context.Context, in *SetValueRequest) (*SetValueResponse, error) {
	return SetValueHandler(ms.App, in)
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
