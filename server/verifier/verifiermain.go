package verifier

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/astromechio/astrocache/cache"
	"github.com/astromechio/astrocache/config"
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/send"
	"github.com/astromechio/astrocache/workers"
	"github.com/pkg/errors"
)

// StartVerifier starts a master node
func StartVerifier() {
	app, err := generateConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "StartVerifier failed to generateConfig"))
	}

	logger.LogInfo("bootstrapping astrocache verifier node(" + app.Self.NID + ")\n")

	startWorkers(app)

	loadChain(app)

	router := router(app)

	addrParts := strings.Split(app.Self.Address, ":")
	port := addrParts[len(addrParts)-1]

	logger.LogInfo(fmt.Sprintf("starting astrocache verifier node server on port %s\n", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatal(err)
	}
}

func startWorkers(app *config.App) {
	go workers.ReserveWorker(app)
	go workers.ProposeWorker(app)
	go workers.CommitWorker(app)
	go workers.ActionWorker(app)
	go workers.DistributeWorker(app)
}

func generateConfig() (*config.App, error) {
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

	chain := blockchain.EmptyChain()

	app := &config.App{
		Self:     node,
		KeySet:   keySet,
		Chain:    chain,
		Cache:    cache.EmptyCache(),
		NodeList: &config.NodeList{},
	}

	newNode, err := send.JoinNetwork(app, masterAddr, joinCode)
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

	app.KeySet.GlobalKey = globalKey

	masterKeyPair, err := acrypto.KeyPairFromPubKeyJSON(newNode.Master.PubKey)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to KeyPairFromPubKeyJSON")
	}

	app.KeySet.AddKeyPair(masterKeyPair)
	app.NodeList.Master = newNode.Master

	if newNode.IsPrimary {
		logger.LogInfo("acting as primary verifier")
		app.NodeList.Master.ParentNID = app.Self.NID
	}

	logger.LogInfo("joined network successfully")

	return app, nil
}

func loadChain(app *config.App) {
	blocks, err := send.GetEntireChain(app.NodeList.Master)
	if err != nil {
		log.Fatal(errors.Wrap(err, "loadChain failed to GetEntireChain, dying now..."))
	}

	if err := app.Chain.LoadFromBlocks(blocks); err != nil {
		log.Fatal(errors.Wrap(err, "loadChain failed to LoadFromBlocks, dying now..."))
	}
}
