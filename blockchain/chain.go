package blockchain

import acrypto "github.com/astromechio/astrocache/crypto"

// Chain represents a blockchain
type Chain struct {
	Blocks []*Block
}

// GenesisBlockID and others are block related consts
const (
	genesisBlockID   = "iamthegenesisbutnottheterminator"
	genesisBlockData = "thegenesisofallblocks"
)

// BrandNewChain creates a fresh chain using the master keyPair
func BrandNewChain(masterKeyPair *acrypto.KeyPair) {

}

func genesisBlock(key *acrypto.KeyPair) (*Block, error) {
	encData, err := key.Encrypt([]byte(genesisBlockData))
	if err != nil {
		return nil, err
	}

	sig, err := key.Sign(encData.Data)
	if err != nil {
		return nil, err
	}

	genesis := &Block{
		ID:        genesisBlockID,
		Data:      encData,
		Signature: sig,
		PrevID:    "",
	}

	return genesis, nil
}
