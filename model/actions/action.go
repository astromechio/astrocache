package actions

import (
	"encoding/json"
	"fmt"

	"github.com/astromechio/astrocache/config"
	"github.com/pkg/errors"
)

// Action defines structs that can be called actions
type Action interface {
	ActionType() string
	JSON() []byte
	Execute(*config.App) error
}

// ActionTypeNodeAdded and others represent different types of actions
const (
	ActionTypeNodeAdded = "astro.action.nodeadded"
)

// UnmarshalAction unmarshals an action from JSON
func UnmarshalAction(actionJSON []byte, actionType string) (Action, error) {
	if actionType == ActionTypeNodeAdded {
		action := &NodeAdded{}
		if err := json.Unmarshal(actionJSON, action); err != nil {
			return nil, errors.Wrap(err, "UnmarshalAction failed to Unmarshal")
		}

		return action, nil
	}

	return nil, fmt.Errorf("UnmarshalAction unable to unmarshal: unknown action type %s", actionType)
}
