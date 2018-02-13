package blockchain

import (
	"fmt"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
	"github.com/pkg/errors"
)

// Chain represents a blockchain
type Chain struct {
	Blocks []*Block
}

// GenesisBlockID and others are block related consts
const (
	genesisBlockID = "iamthegenesisbutnottheterminator"
)

// AddBlock verifies and adds a block to the chain
func (c *Chain) AddBlock(block *Block, keySet *acrypto.KeySet) error {
	prev := c.Blocks[len(c.Blocks)-1]

	if err := block.Verify(keySet, prev); err != nil {
		return errors.Wrap(err, "AddBlock failed to block.Verify")
	}

	c.Blocks = append(c.Blocks, block)

	return nil
}

// BrandNewChain creates a fresh chain using the master keyPair
func BrandNewChain(masterKeyPair *acrypto.KeyPair, globalKey *acrypto.SymKey, node *model.Node) (*Chain, error) {
	if masterKeyPair.KID != acrypto.MasterKeyPairKID {
		return nil, fmt.Errorf("attempted to create new chain with non-master keyPair")
	}

	if globalKey.KID != acrypto.GlobalSymKeyKID {
		return nil, fmt.Errorf("attempted to create new chain with non-global symKey")
	}

	nodeKeyPair, err := node.KeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "BrandNewChain failed to node.KeyPair")
	}

	if nodeKeyPair.KID != masterKeyPair.KID || node.Type != model.NodeTypeMaster {
		return nil, fmt.Errorf("attempted to create a new chain with a non-master node")
	}

	globalKeyJSON := globalKey.JSON()

	encGlobalKey, err := masterKeyPair.Encrypt(globalKeyJSON)
	if err != nil {
		return nil, errors.Wrap(err, "BrandNewChain failed to masterKeyPair.Encrypt")
	}

	nodeAddedAction := model.NewNodeAddedAction(node, encGlobalKey)

	genesis, err := NewBlockWithAction(masterKeyPair, globalKey, nodeAddedAction, nil)
	if err != nil {
		return nil, errors.Wrap(err, "BrandNewChain failed to NewBlockWithAction")
	}

	chain := &Chain{
		Blocks: []*Block{genesis},
	}

	return chain, nil
}
