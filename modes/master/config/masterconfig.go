package config

import (
	"math/rand"

	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/modes"
)

// Config defines the configuration for a master node
type Config struct {
	*modes.BaseConfig
	Nodes    *NodeList
	JoinCode string
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
	nl.Workers = append(nl.Workers, worker)
}
