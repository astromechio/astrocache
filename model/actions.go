package model

import (
	acrypto "github.com/astromechio/astrocache/crypto"
)

// NodeAddedAction is a block value representing a new node in the network
type NodeAddedAction struct {
	ID       string                      `json:"id"`
	PubKey   *acrypto.SerializablePubKey `json:"pubKey"`
	SymKey   *acrypto.Message            `json:"encSymKey"`
	NodeType string                      `json:"nodeType"`
}

// NodeTypeMaster and others represent node types
const (
	NodeTypeMaster   = "astromaster"
	NodeTypeVerifier = "astroverifier"
	NodeTypeWorker   = "astroworker"
)
