package verifier

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/astromechio/astrocache/execute"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/modes"
	"github.com/astromechio/astrocache/modes/verifier/config"
	"github.com/pkg/errors"
)

// StartVerifier starts a master node
func StartVerifier() {
	logger.LogInfo("bootstrapping astrocache verifier node")

	config, err := generateConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "StartVerifier failed to generateConfig"))
	}

	router := router(config)

	addrParts := strings.Split(config.Self.Address, ":")
	port := addrParts[len(addrParts)-1]
	fmt.Println(port)

	logger.LogInfo(fmt.Sprintf("starting astrocache verifier node server on port %s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
}

func generateConfig() (*config.Config, error) {
	if len(os.Args) < 3 {
		return nil, errors.New("missing argument: address")
	}

	if len(os.Args) < 4 {
		return nil, errors.New("missing argument: master node address")
	}

	if len(os.Args) < 5 {
		return nil, errors.New("missing argument: join code")
	}

	address := os.Args[2]
	if strings.Index(address, ":") < 0 {
		return nil, errors.New("address does not contain port value")
	}

	masterAddr := os.Args[3]
	joinCode := os.Args[4]

	keyPair, err := acrypto.GenerateNewKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to GenerateMasterKeyPair")
	}

	node := model.NewNode(address, model.NodeTypeVerifier, keyPair)

	keySet := &acrypto.KeySet{
		KeyPair: keyPair,
	}

	chain := &blockchain.Chain{
		Blocks:  []*blockchain.Block{},
		Pending: []*blockchain.Block{},
	}

	config := &config.Config{
		BaseConfig: &modes.BaseConfig{
			Self:   node,
			KeySet: keySet,
			Chain:  chain,
		},
		Workers: []*model.Node{},
	}

	encGlobalKey, err := joinNetwork(masterAddr, joinCode, config)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to joinNetwork")
	}

	globalKeyJSON, err := keyPair.Decrypt(encGlobalKey)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to Decrypt")
	}

	globalKey, err := acrypto.SymKeyFromJSON(globalKeyJSON)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to SymKeyFromJSON")
	}

	config.KeySet.GlobalKey = globalKey

	return config, nil
}

func joinNetwork(masterAddr, joinCode string, config *config.Config) (*acrypto.Message, error) {
	logger.LogInfo(fmt.Sprintf("joining network with master node at address %s", masterAddr))

	newNode, err := execute.JoinNetwork(config.BaseConfig, masterAddr, joinCode)
	if err != nil {
		return nil, errors.Wrap(err, "joinNetwork failed to JoinNetwork")
	}

	return newNode.EncGlobalKey, nil
}
