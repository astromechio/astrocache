package modes

// for the record, I hate this package name but w/e

import (
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
)

// BaseConfig defines the configuration for a master node
type BaseConfig struct {
	Self   *model.Node
	KeySet *acrypto.KeySet
}
