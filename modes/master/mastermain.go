package master

import (
	"log"
	"net/http"
	"os"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/modes"
	"github.com/pkg/errors"
)

// StartMaster starts a master node
func StartMaster() {
	logger.LogInfo("bootstrapping astrocache master node")

	config, err := generateConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "StartMaster failed to generateConfig"))
	}

	router := router(config)

	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}

func generateConfig() (*Config, error) {
	address := os.Args[2]

	keyPair, err := acrypto.GenerateMasterKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to GenerateMasterKeyPair")
	}

	globalKey, err := acrypto.GenerateGlobalSymKey()
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to GenerateGlobalSymKey")
	}

	node := model.NewNode(address, model.NodeTypeMaster, keyPair)

	keySet := &acrypto.KeySet{
		KeyPair:   keyPair,
		GlobalKey: globalKey,
	}

	config := &Config{
		BaseConfig: &modes.BaseConfig{
			Self:   node,
			KeySet: keySet,
		},
		Nodes: &NodeList{},
	}

	return config, nil
}
