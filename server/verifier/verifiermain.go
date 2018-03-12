package verifier

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/send"
	"github.com/astromechio/astrocache/server"
	"github.com/astromechio/astrocache/workers"
	"github.com/pkg/errors"
)

// StartVerifier starts a master node
func StartVerifier() {
	logger.LogInfo("bootstrapping astrocache verifier node\n")

	config, err := generateConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "StartVerifier failed to generateConfig"))
	}

	go workers.StartChainWorker(config)
	go loadChain(config)

	router := router(config)

	addrParts := strings.Split(config.Self.Address, ":")
	port := addrParts[len(addrParts)-1]

	logger.LogInfo(fmt.Sprintf("starting astrocache verifier node server on port %s", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
}

func generateConfig() (*server.Config, error) {
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
	masterNode := &model.Node{
		Address: fmt.Sprintf("%s:3000", masterAddr),
	}

	joinCode := os.Args[4]

	keyPair, err := acrypto.GenerateNewKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to GenerateMasterKeyPair")
	}

	node := model.NewNode(address, model.NodeTypeVerifier, keyPair)

	keySet := &acrypto.KeySet{
		KeyPair: keyPair,
	}

	chain := blockchain.EmptyChain()

	config := &server.Config{
		Self:   node,
		KeySet: keySet,
		Chain:  chain,
		NodeList: &server.NodeList{
			Master: masterNode,
		},
	}

	newNode, err := send.JoinNetwork(config, masterAddr, joinCode)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to JoinNetwork")
	}

	globalKeyJSON, err := keyPair.Decrypt(newNode.EncGlobalKey)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to Decrypt")
	}

	globalKey, err := acrypto.SymKeyFromJSON(globalKeyJSON)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to SymKeyFromJSON")
	}

	config.KeySet.GlobalKey = globalKey

	masterKeyPair, err := acrypto.KeyPairFromPubKeyJSON(newNode.MasterPubKeyJSON)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to KeyPairFromPubKeyJSON")
	}

	config.KeySet.AddKeyPair(masterKeyPair)

	logger.LogInfo("joined network successfully")

	return config, nil
}

func loadChain(config *server.Config) {
	blocks, err := send.GetEntireChain(config.NodeList.Master)
	if err != nil {
		log.Fatal(errors.Wrap(err, "loadChain failed to GetEntireChain, dying now..."))
	}

	if err := config.Chain.LoadFromBlocks(blocks); err != nil {
		log.Fatal(errors.Wrap(err, "loadChain failed to LoadFromBlocks, dying now..."))
	}
}
