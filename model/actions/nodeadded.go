package actions

import (
	"encoding/json"
	"fmt"

	"github.com/astromechio/astrocache/config"
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/pkg/errors"
)

// NodeAdded is a block value representing a new node in the network
// GlobalKey is the global key encrypted with the new node's pubKey
type NodeAdded struct {
	Node         *model.Node      `json:"node"`
	EncGlobalKey *acrypto.Message `json:"encGlobalKey"`
}

// NewNodeAdded creates a new NodeAdded
func NewNodeAdded(node *model.Node, globalKey *acrypto.Message) *NodeAdded {
	return &NodeAdded{
		Node:         node,
		EncGlobalKey: globalKey,
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

// Execute adds the node to the node list
func (na *NodeAdded) Execute(app *config.App) error {
	logger.LogInfo("Adding node with NID " + na.Node.NID)

	if na.Node.NID == app.Self.NID {
		logger.LogInfo("NodeAdded.Execute tried to add self, skipping...")
		return nil
	}

	pubKey, err := acrypto.KeyPairFromPubKeyJSON(na.Node.PubKey)
	if err != nil {
		return errors.Wrap(err, "NodeAdded.Execute failed to KeyPairFromPubKeyJSON")
	}

	app.KeySet.AddKeyPair(pubKey)

	// worker nodes only need to know about the master and their verifier, so skip the rest
	if app.Self.Type == model.NodeTypeWorker {
		return nil
	}

	if na.Node.Type == model.NodeTypeVerifier {
		app.NodeList.AddVerifier(na.Node)
	} else if na.Node.Type == model.NodeTypeWorker {
		app.NodeList.AddWorker(na.Node)
	} else if na.Node.Type == model.NodeTypeMaster {
		if app.NodeList.Master == nil {
			app.NodeList.Master = na.Node
		} else {
			logger.LogWarn("NodeAdded.Execute tried to set a master node when one already exists, skipping...")
		}
	} else {
		return fmt.Errorf("NodeAdded.Execute tried to add node with unknown type %q", na.Node.Type)
	}

	return nil
}
