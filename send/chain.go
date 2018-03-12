package send

import (
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/transport"
	"github.com/pkg/errors"
)

// GetEntireChain requests the entire chain from the master node
func GetEntireChain(masterNode *model.Node) ([]*blockchain.Block, error) {
	url := transport.URLFromAddressAndPath(masterNode.Address, "v1/master/chain")

	blocks := []*blockchain.Block{}
	if err := transport.Get(url, &blocks); err != nil {
		return nil, errors.Wrap(err, "GetEntireChain failed to Get")
	}

	return blocks, nil
}
