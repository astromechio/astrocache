package verifier

import (
	"log"
	"net/http"
	"os"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/modes"
	"github.com/pkg/errors"
)

// StartVerifier starts a master node
func StartVerifier() {
	logger.LogInfo("bootstrapping astrocache verifier node")

	config, err := generateConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "StartVerifier failed to generateConfig"))
	}

	chain := &blockchain.Chain{
		Blocks:  []*blockchain.Block{},
		Pending: []*blockchain.Block{},
	}

	router := router(config, chain)

	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}

func generateConfig() (*Config, error) {
	address := os.Args[2]

	keyPair, err := acrypto.GenerateNewKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to GenerateMasterKeyPair")
	}

	node := model.NewNode(address, model.NodeTypeVerifier, keyPair)

	keySet := &acrypto.KeySet{
		KeyPair: keyPair,
	}

	config := &Config{
		BaseConfig: &modes.BaseConfig{
			Self:   node,
			KeySet: keySet,
		},
		Workers: []*model.Node{},
	}

	return config, nil
}
