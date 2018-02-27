package master

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/modes"
	"github.com/astromechio/astrocache/modes/master/config"
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

	logger.LogInfo(fmt.Sprintf("to join the network, run `astrocache [worker|verifier] [node address] %s %s`\n", config.Self.Address, config.JoinCode))

	logger.LogInfo("starting astrocache master node server on port 3000")
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}

func generateConfig() (*config.Config, error) {
	if len(os.Args) < 3 {
		return nil, errors.New("missing argument: address")
	}

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

	chain, err := blockchain.BrandNewChain(keyPair, globalKey, node)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to BrandNewChain")
	}

	config := &config.Config{
		BaseConfig: &modes.BaseConfig{
			Self:   node,
			KeySet: keySet,
			Chain:  chain,
		},
		Nodes:    &config.NodeList{},
		JoinCode: generateJoinCode(),
	}

	return config, nil
}

func generateJoinCode() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return acrypto.Base64URLEncode(bytes)
}
