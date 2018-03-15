package blockchain

import (
	"fmt"

	"github.com/astromechio/astrocache/logger"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/pkg/errors"
)

// Chain represents a blockchain
type Chain struct {
	Blocks   []*Block
	Proposed *Block

	ProposeChan chan (*NewBlockJob) // ProposeChan is used by proposeworker as the synchronization method for proposing blocks
	VerifyChan  chan (*NewBlockJob) // Check is used by proposeworker as the synchronization method for processing incoming blocks
	CommitChan  chan (*NewBlockJob) // CommitChan is used by commitworker as the synchronization method for committing blocks

	ProposedChan  chan (*Block) // ProposedChan is used when a goroutine needs to know the next time a block is proposed.
	CommittedChan chan (*Block) // CommittedChan is used when a goroutine needs to know the next time a block is committed.

	ActionChan     chan (*Block) // ActionChan decrypts blocks and applies actions
	DistributeChan chan (*Block) // DistributeChan loads blocks needed to be distributed to workers
}

// NewBlockJob represents the intent to add a new block
type NewBlockJob struct {
	Block        *Block
	ProposingNID string
	ResultChan   chan (error)
	Check        bool
}

// AddNewBlock checks and then sets the proposed block
func (c *Chain) AddNewBlock(block *Block, propNID string) chan error {
	errChan := make(chan error, 1)

	job := &NewBlockJob{
		Block:        block,
		ProposingNID: propNID,
		ResultChan:   errChan,
		Check:        true,
	}

	c.sendProposeJob(job)

	return errChan
}

// VerifyProposedBlock checks and then sets the proposed block
func (c *Chain) VerifyProposedBlock(block *Block, propNID string) chan error {
	errChan := make(chan error, 1)

	job := &NewBlockJob{
		Block:        block,
		ProposingNID: propNID,
		ResultChan:   errChan,
		Check:        true,
	}

	c.sendVerifyJob(job)

	return errChan
}

// VerifyBlockUnchecked checks and then sets the proposed block
func (c *Chain) VerifyBlockUnchecked(block *Block) chan error {
	errChan := make(chan error, 1)

	job := &NewBlockJob{
		Block:      block,
		ResultChan: errChan,
		Check:      false,
	}

	c.sendVerifyJob(job)

	return errChan
}

func (c *Chain) sendProposeJob(job *NewBlockJob) {
	c.ProposeChan <- job
}

func (c *Chain) sendVerifyJob(job *NewBlockJob) {
	c.VerifyChan <- job
}

// LoadFromBlocks loads a chain from a block array
func (c *Chain) LoadFromBlocks(blocks []*Block) error {
	if len(c.Blocks) > 0 {
		return fmt.Errorf("LoadFromBlocks attempted to load chain with %d existing blocks", len(c.Blocks))
	}

	for i := range blocks {
		errChan := c.VerifyBlockUnchecked(blocks[i])
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// GetNextProposed blocks until the next block is proposed
func (c *Chain) GetNextProposed() chan *Block {
	logger.LogInfo("Waiting for next proposed block")

	if c.ProposedChan != nil {
		return nil
	}

	notifChan := make(chan *Block, 1)
	c.ProposedChan = notifChan

	return notifChan
}

// GetNextCommitted blocks until the next block is proposed
func (c *Chain) GetNextCommitted() chan *Block {
	logger.LogInfo("Waiting for next committed block")

	if c.CommittedChan != nil {
		return nil
	}

	notifChan := make(chan *Block, 1)
	c.CommittedChan = notifChan

	return notifChan
}

// HasProposedOrCommittedBlock checks if a block is proposed or previously committed
func (c *Chain) HasProposedOrCommittedBlock(block *Block) bool {
	if c.Proposed != nil && c.Proposed.IsSameAsBlock(block) {
		return true
	}

	// check if the block being checked is the next to be committed
	last := c.LastBlock()
	if last != nil {
		hash, err := last.Hash()
		if err != nil {
			return false
		}

		lastHashString := acrypto.Base64URLEncode(hash)
		if lastHashString == block.ID {
			return true
		}
	}

	// if that fails, work backwards in the chain to see if any of them match
	for i := len(c.Blocks) - 1; i > 0; i-- {
		if c.Blocks[i].IsSameAsBlock(block) {
			return true
		}
	}

	return false
}

// BlocksAfterID returns all the committed blocks after id
func (c *Chain) BlocksAfterID(id string) []*Block {
	for i := len(c.Blocks) - 1; i > 0; i-- {
		if c.Blocks[i].ID == id {
			return c.Blocks[i+1:]
		}
	}

	return nil
}

// EmptyChain creates an enpty chain
func EmptyChain() *Chain {
	chain := &Chain{
		Blocks:      []*Block{},
		ProposeChan: make(chan *NewBlockJob),
		VerifyChan:  make(chan *NewBlockJob, 20),
		CommitChan:  make(chan *NewBlockJob, 2),
		// ProposedChan:   make(chan *Block),
		CommittedChan:  make(chan *Block, 1), // this is buffered because only the ProposeWorker cares about it
		ActionChan:     make(chan *Block),
		DistributeChan: make(chan *Block),
	}

	return chain
}

// BrandNewChain creates a fresh chain using the master keyPair
func BrandNewChain(masterKeyPair *acrypto.KeyPair, globalKey *acrypto.SymKey, blockData []byte, actionType string) (*Chain, error) {
	if masterKeyPair.KID != acrypto.MasterKeyPairKID {
		return nil, fmt.Errorf("attempted to create new chain with non-master keyPair")
	}

	if globalKey.KID != acrypto.GlobalSymKeyKID {
		return nil, fmt.Errorf("attempted to create new chain with non-global symKey")
	}

	genesis, err := genesisBlockWithData(globalKey, blockData, actionType)
	if err != nil {
		return nil, errors.Wrap(err, "BrandNewChain failed to NewBlockWithAction")
	}

	if err := genesis.PrepareForCommit(masterKeyPair, nil); err != nil {
		return nil, errors.Wrap(err, "BrandNewChain failed to PrepareForCommit")
	}

	chain := EmptyChain()

	// if this fails in the worker, we'll have to catch it and fatal
	chain.Blocks = append(chain.Blocks, genesis)

	return chain, nil
}

// LastBlock returns the last block in the chain
func (c *Chain) LastBlock() *Block {
	if len(c.Blocks) == 0 {
		return nil
	}

	return c.Blocks[len(c.Blocks)-1]
}
