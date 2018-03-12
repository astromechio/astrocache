package server

import (
	"math/rand"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
)

// ConfigJoinCodeKey and others are config related constants
const (
	ConfigJoinCodeKey = "astro.master.joincode"
)

// Config defines the configuration for a node
type Config struct {
	Self     *model.Node
	KeySet   *acrypto.KeySet
	Chain    *blockchain.Chain
	NodeList *NodeList
	Values   map[string]string
}

// SetValueForKey sets a value for a key
func (c *Config) SetValueForKey(val, key string) {
	if c.Values == nil {
		c.Values = make(map[string]string)
	}

	c.Values[key] = val
}

// ValueForKey retreives a value for a particular key
func (c *Config) ValueForKey(key string) string {
	return c.Values[key]
}

// NodeList defines the nodes a master looks after
type NodeList struct {
	Verifiers []*model.Node
	Workers   []*model.Node
}

// WorkersForVerifierWithNID returns the worker nodes assigned to a verifier node with NID
func (nl *NodeList) WorkersForVerifierWithNID(nid string) []*model.Node {
	workers := []*model.Node{}

	for i, w := range nl.Workers {
		if w.ParentNID == nid {
			workers = append(workers, nl.Workers[i])
		}
	}

	return workers
}

// AddVerifier adds a verifier to the nodeList
func (nl *NodeList) AddVerifier(verifier *model.Node) {
	if nl.Verifiers == nil {
		nl.Verifiers = []*model.Node{}
	}

	nl.Verifiers = append(nl.Verifiers, verifier)
}

// RandomVerifier returns a random verifier node from the NodeList
func (nl *NodeList) RandomVerifier() *model.Node {
	if len(nl.Verifiers) == 0 {
		return nil
	} else if len(nl.Verifiers) == 1 {
		return nl.Verifiers[0]
	}

	index := rand.Intn(len(nl.Verifiers))

	return nl.Verifiers[index]
}

// AddWorker adds a worker to the nodeList
func (nl *NodeList) AddWorker(worker *model.Node) {
	if nl.Workers == nil {
		nl.Workers = []*model.Node{}
	}

	nl.Workers = append(nl.Workers, worker)
}
