package model

import (
	"encoding/json"

	acrypto "github.com/astromechio/astrocache/crypto"
)

// NodeAddedAction is a block value representing a new node in the network
// GlobalKey is the global key encrypted with the new node's pubKey
type NodeAddedAction struct {
	Node      *Node            `json:"node"`
	GlobalKey *acrypto.Message `json:"globalKey"`
}

// NewNodeAddedAction creates a new NodeAddedAction
func NewNodeAddedAction(node *Node, globalKey *acrypto.Message) *NodeAddedAction {
	return &NodeAddedAction{
		Node:      node,
		GlobalKey: globalKey,
	}
}

// ActionType defines this action's type
func (na *NodeAddedAction) ActionType() string {
	return ActionTypeNodeAdded
}

// JSON returns json for the action
func (na *NodeAddedAction) JSON() []byte {
	naJSON, _ := json.Marshal(na)

	return naJSON
}
