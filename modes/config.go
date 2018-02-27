package modes

// for the record, I hate this package name but w/e

import (
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
)

// BaseConfig defines the configuration for a master node
type BaseConfig struct {
	Self   *model.Node
	KeySet *acrypto.KeySet
	Chain  *blockchain.Chain
}
