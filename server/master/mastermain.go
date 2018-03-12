package master

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/astromechio/astrocache/model/actions"

	"github.com/astromechio/astrocache/config"
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/workers"
	"github.com/pkg/errors"
)

// StartMaster starts a master node
func StartMaster() {
	app, err := generateConfig()
	if err != nil {
		log.Fatal(errors.Wrap(err, "StartMaster failed to generateConfig"))
	}

	logger.LogInfo("bootstrapping astrocache master node (" + app.Self.NID + ")\n")

	joinCode := app.ValueForKey(config.AppJoinCodeKey)
	logger.LogInfo(fmt.Sprintf("to join the network, run `astrocache [worker|verifier] [node address] %s %s`\n", app.Self.Address, joinCode))

	startWorkers(app)

	router := router(app)

	logger.LogInfo("starting astrocache master node server on port 3000\n")
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}

func startWorkers(app *config.App) {
	go workers.StartActionWorker(app)
	go workers.StartChainWorker(app)
}

func generateConfig() (*config.App, error) {
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

	globalKeyJSON := globalKey.JSON()

	encGlobalKey, err := keyPair.Encrypt(globalKeyJSON)
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to Encrypt")
	}

	nodeAddedAction := actions.NewNodeAdded(node, encGlobalKey)
	actionJSON := nodeAddedAction.JSON()

	chain, err := blockchain.BrandNewChain(keyPair, globalKey, actionJSON, nodeAddedAction.ActionType())
	if err != nil {
		return nil, errors.Wrap(err, "generateConfig failed to BrandNewChain")
	}

	app := &config.App{
		Self:     node,
		KeySet:   keySet,
		Chain:    chain,
		NodeList: &config.NodeList{},
	}

	joinCode := generateJoinCode()
	app.SetValueForKey(joinCode, config.AppJoinCodeKey)

	return app, nil
}

func generateJoinCode() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return acrypto.Base64URLEncode(bytes)
}
