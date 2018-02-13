package actions

import (
	"encoding/json"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
)

// NodeAdded is a block value representing a new node in the network
// GlobalKey is the global key encrypted with the new node's pubKey
type NodeAdded struct {
	Node      *model.Node      `json:"node"`
	GlobalKey *acrypto.Message `json:"globalKey"`
}

// NewNodeAdded creates a new NodeAdded
func NewNodeAdded(node *model.Node, globalKey *acrypto.Message) *NodeAdded {
	return &NodeAdded{
		Node:      node,
		GlobalKey: globalKey,
	}
}

// ActionType defines this action's type
func (na *NodeAdded) ActionType() string {
	return ActionTypeNodeAdded
}

// JSON returns json for the action
func (na *NodeAdded) JSON() []byte {
	naJSON, _ := json.Marshal(na)

	return naJSON
}
