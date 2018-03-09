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
	Blocks  []*Block
	Pending []*Block
}

// AddPendingBlock adds a block to the pending list and returns the potential prevBlock
func (c *Chain) AddPendingBlock(block *Block, keySet *acrypto.KeySet) error {
	prevBlock := c.LastBlock()

	if prevBlock.ID != block.ID {
		return errors.New("AddPendingBlock failed to add block: block.PrevID did not match prevBlock.ID")
	}

	c.Pending = append(c.Pending, block)

	return nil
}

// CommitBlockWithTempID adds prepares a block for committing and then commits it
// If newSig is nil, we are committing this block and therefore must prepare it
func (c *Chain) CommitBlockWithID(tempID, prevID string, newSig *acrypto.Signature, keySet *acrypto.KeySet) error {
	var prevBlock *Block
	var newBlock *Block
	isNextBlock := false

	for i, b := range c.Pending {
		if b.ID == tempID {
			newBlock = c.Pending[i]

			if i == 0 {
				prevBlock = c.Blocks[len(c.Blocks)-1]
				isNextBlock = true
			} else {
				prevBlock = c.Pending[i-1]
			}

			break
		}
	}

	if newBlock == nil {
		return fmt.Errorf("CommitBlockWithTempID unable to find pending block with tempID %s", tempID)
	}

	// if we are the committer
	if newSig == nil {
		newBlock.prepareForCommit(keySet.KeyPair, prevBlock)
	} else {
		// if we are not the committer
		prevHash, err := prevBlock.Hash()
		if err != nil {
			return errors.Wrap(err, "CommitBlockWithTempID failed to prevBlock.Hash")
		}

		newBlock.ID = acrypto.Base64URLEncode(prevHash)
		newBlock.Signature = newSig
		newBlock.PrevID = prevID

		if err := newBlock.Verify(keySet, prevBlock); err != nil {
			return errors.Wrap(err, "CommitBlockWithTempID failed to block.Verify")
		}
	}

	logger.LogInfo(fmt.Sprintf("*** Committing bock with ID %s ***", newBlock.ID))
	c.Blocks = append(c.Blocks, newBlock)

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

	nodeAddedAction := actions.NewNodeAdded(node, encGlobalKey)

	genesis, err := genesisBlockWithAction(globalKey, nodeAddedAction)
	if err != nil {
		return nil, errors.Wrap(err, "BrandNewChain failed to NewBlockWithAction")
	}

	chain := &Chain{
		Blocks: []*Block{genesis},
	}

	return chain, nil
}

// LastBlock returns the last [pending | committed] block in the chain
func (c *Chain) LastBlock() *Block {
	var lastBlock *Block

	if len(c.Pending) > 0 {
		lastBlock = c.Pending[len(c.Pending)-1]
	} else {
		lastBlock = c.Blocks[len(c.Blocks)-1]
	}

	return lastBlock
}

// LastCommittedBlock returns the last block that has been committed
func (c *Chain) LastCommittedBlock() *Block {
	return c.Blocks[len(c.Blocks)-1]
}
