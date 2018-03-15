package actions

import (
	"encoding/json"
	"fmt"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
)

// SetValue is a block value representing a new node in the network
// GlobalKey is the global key encrypted with the new node's pubKey
type SetValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewSetValue creates a nSetValue
func NewSetValue(key, value string) *SetValue {
	return &SetValue{
		Key:   key,
		Value: value,
	}
}

// ActionType defines this action's type
func (sv *SetValue) ActionType() string {
	return ActionTypeSetValue
}

// JSON returns json for the action
func (sv *SetValue) JSON() []byte {
	naJSON, _ := json.Marshal(sv)

	return naJSON
}

// Execute adds the node to the node list
func (sv *SetValue) Execute(app *config.App) error {
	if app.Self.Type == model.NodeTypeWorker {
		logger.LogInfo(fmt.Sprintf("Setting value %q for key %q", sv.Value, sv.Key))

		app.Cache.SetValueForKey(sv.Value, sv.Key)
	}

	return nil
}
