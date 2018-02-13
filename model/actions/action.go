package actions

// Action defines structs that can be called actions
type Action interface {
	ActionType() string
	JSON() []byte
}

// ActionTypeNodeAdded and others represent different types of actions
const (
	ActionTypeNodeAdded = "astro.action.nodeadded"
)
