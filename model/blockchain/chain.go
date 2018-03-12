package blockchain

import (
	"fmt"

	"github.com/astromechio/astrocache/logger"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/actions"
	"github.com/pkg/errors"
)

// Chain represents a blockchain
type Chain struct {
	Blocks     []*Block
	Proposed   *Block
	WorkerChan chan (*NewBlockJob) // WorkerChan is used by verifier_chainworker as the synchronization method for mining and verifying new blocks
}

// NewBlockJob represents the intent to add a new block
type NewBlockJob struct {
	Block      *Block
	ResultChan chan (error)
}

// AddNewBlock adds a block to the pending list and returns the potential prevBlock
func (c *Chain) AddNewBlock(block *Block) chan (error) {
	errChan := make(chan error, 1)

	job := &NewBlockJob{
		Block:      block,
		ResultChan: errChan,
	}

	c.WorkerChan <- job
	return errChan
}

// SetProposedBlock checks and then sets the proposed block
func (c *Chain) SetProposedBlock(block *Block) error {
	prevBlock := c.LastBlock()

	if prevBlock == nil {
		if block.ID != genesisBlockID {
			return fmt.Errorf("SetProposedBlock tried to propose block with nil prevBlock and non-genesis ID %q", block.ID)
		}
	} else {
		if prevBlock.ID != block.PrevID {
			return fmt.Errorf("AddPendingBlock failed to add block: block.PrevID (%s) did not match prevBlock.ID (%s)", block.PrevID, prevBlock.ID)
		}
	}

	c.Proposed = block

	return nil
}

// CommitProposedBlock verifies and commits a block
func (c *Chain) CommitProposedBlock(keySet *acrypto.KeySet) error {
	prevBlock := c.LastBlock()

	// Verify handles the genesis case
	if err := c.Proposed.Verify(keySet, prevBlock); err != nil {
		return errors.Wrap(err, "CommitProposedBlock failed to block.Verify")
	}

	logger.LogInfo(fmt.Sprintf("*** Committing bock with ID %s ***", c.Proposed.ID))

	c.Blocks = append(c.Blocks, c.Proposed)
	c.Proposed = nil

	return nil
}

// LoadFromBlocks loads a chain from a block array
func (c *Chain) LoadFromBlocks(blocks []*Block) error {
	if len(c.Blocks) > 0 {
		return fmt.Errorf("LoadFromBlocks attempted to load chain with %d existing blocks", len(c.Blocks))
	}

	for i := range blocks {
		errChan := c.AddNewBlock(blocks[i])
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// EmptyChain creates an enpty chain
func EmptyChain() *Chain {
	chain := &Chain{
		Blocks:     []*Block{},
		WorkerChan: make(chan *NewBlockJob, 2),
	}

	return chain
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
		return nil, fmt.Errorf("BrandNewChain attempted to create a new chain with a non-master node or keypair")
	}

	globalKeyJSON := globalKey.JSON()

	encGlobalKey, err := masterKeyPair.Encrypt(globalKeyJSON)
	if err != nil {
		return nil, errors.Wrap(err, "BrandNewChain failed to masterKeyPair.Encrypt")
	}

	nodeAddedAction := actions.NewNodeAdded(node, encGlobalKey)

	genesis, err := genesisBlockWithAction(globalKey, nodeAddedAction)
	if err != nil {
		return nil, errors.Wrap(err, "BrandNewChain failed to NewBlockWithAction")
	}

	chain := EmptyChain()

	// if this fails in the worker, we'll have to catch it and fatal
	chain.AddNewBlock(genesis)

	return chain, nil
}

// LastBlock returns the last block in the chain
func (c *Chain) LastBlock() *Block {
	if len(c.Blocks) == 0 {
		return nil
	}

	return c.Blocks[len(c.Blocks)-1]
}
