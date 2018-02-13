package verifier

import (
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/modes"
)

// Config defines the configuration for a master node
type Config struct {
	*modes.BaseConfig
	Workers []*model.Node
}

// AddWorker adds a worker to the nodeList
func (c *Config) AddWorker(worker *model.Node) {
	c.Workers = append(c.Workers, worker)
}
