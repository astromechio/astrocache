package model

import (
	"crypto/rand"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// NodeTypeMaster and others represent node types
const (
	NodeTypeMaster   = "astro.cache.master"
	NodeTypeVerifier = "astro.cache.verifier"
	NodeTypeWorker   = "astro.cache.worker"
)

// Node defines a node in the network
type Node struct {
	NID       string           `json:"nid"`
	Address   string           `json:"address"`
	Type      string           `json:"type"`
	PubKey    []byte           `json:"pubKey"`
	ParentNID string           `json:"parentNid,omitempty"`
	Conn      *grpc.ClientConn `json:"-"`
}

// NewNode creates a new node
func NewNode(addr, nodeType string, keyPair *acrypto.KeyPair) *Node {
	nid := generateNewNID()

	pubKeyJSON := keyPair.PubKeyJSON()

	return &Node{
		NID:     nid,
		Address: addr,
		Type:    nodeType,
		PubKey:  pubKeyJSON,
	}
}

// KeyPair returns the node's pubKey
func (n *Node) KeyPair() (*acrypto.KeyPair, error) {
	return acrypto.KeyPairFromPubKeyJSON(n.PubKey)
}

// Dial configures the node's grpc connection
func (n *Node) Dial() (*grpc.ClientConn, error) {
	if n.Conn != nil {
		return n.Conn, nil
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(n.Address, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "Dial failed to grpc.Dial")
	}

	n.Conn = conn

	return conn, nil
}

func generateNewNID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return acrypto.Base64URLEncode(bytes)
}
