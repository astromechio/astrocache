package service

import (
	fmt "fmt"

	"github.com/astromechio/astrocache/config"
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/actions"
	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/pkg/errors"
)

// AddNodeHandler handles node add rpcs for a master node
func AddNodeHandler(app *config.App, req *NewNodeRequest) (*NewNodeResponse, error) {
	var err error

	node := NodeFromProtoNode(req.Node)

	joinCode := app.ValueForKey(config.AppJoinCodeKey)
	if req.JoinCode != joinCode {
		err = errors.New("AddVerifierNodeHandler failed to verify JoinCode")
		logger.LogError(err)
		return nil, err
	}

	if node.Type == model.NodeTypeVerifier {
		return addVerifierNode(app, node)
	} else if node.Type == model.NodeTypeWorker {
		return addWorkerNode(app, node)
	}

	return nil, fmt.Errorf("AddNodeHandler encountered unsupported node type %s", node.Type)
}

func addVerifierNode(app *config.App, node *model.Node) (*NewNodeResponse, error) {
	var err error

	newNodePubKey, err := acrypto.KeyPairFromPubKeyJSON(node.PubKey)
	if err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to KeyPairFromPubKeyJSON")
		logger.LogError(err)
		return nil, err
	}

	encGlobalKey, err := newNodePubKey.Encrypt(app.KeySet.GlobalKey.JSON())
	if err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to Encrypt")
		logger.LogError(err)
		return nil, err
	}

	nodeAddedAction := actions.NewNodeAdded(node, encGlobalKey)
	actionJSON := nodeAddedAction.JSON()

	block, err := blockchain.NewBlockWithData(app.KeySet.GlobalKey, actionJSON, nodeAddedAction.ActionType())
	if err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to NewBlockWithAction")
		logger.LogError(err)
		return nil, err
	}

	errChan, _ := app.Chain.ReserveBlockID(app.Self.NID)
	if err = <-errChan; err != nil {
		err = errors.Wrap(err, "SetValueHandler failed to ReserveBlockID")
		logger.LogError(err)
		return nil, err
	}

	errChan = app.Chain.AddNewBlock(block, app.Self.NID)
	if err = <-errChan; err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to AddNewBlock")
		logger.LogError(err)
		return nil, err
	}

	encGlobalKeyJSON, err := encGlobalKey.ToJSON()
	if err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to ToJSON")
		logger.LogError(err)
		return nil, err
	}

	protoMaster := ProtoNodeFromNode(app.Self)

	resp := &NewNodeResponse{
		EncGlobalKey: encGlobalKeyJSON,
		Master:       protoMaster,
	}

	// if this is the first verifier, it will be responsible for distributing blocks to us
	resp.IsPrimary = len(app.NodeList.Verifiers) == 0

	return resp, nil
}

func addWorkerNode(app *config.App, node *model.Node) (*NewNodeResponse, error) {
	var err error

	newNodePubKey, err := acrypto.KeyPairFromPubKeyJSON(node.PubKey)
	if err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to KeyPairFromPubKeyJSON")
		logger.LogError(err)
		return nil, err
	}

	encGlobalKey, err := newNodePubKey.Encrypt(app.KeySet.GlobalKey.JSON())
	if err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to Encrypt")
		logger.LogError(err)
		return nil, err
	}

	verifier := app.NodeList.RandomVerifier()
	node.ParentNID = verifier.NID

	nodeAddedAction := actions.NewNodeAdded(node, encGlobalKey)
	actionJSON := nodeAddedAction.JSON()

	block, err := blockchain.NewBlockWithData(app.KeySet.GlobalKey, actionJSON, nodeAddedAction.ActionType())
	if err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to NewBlockWithAction")
		logger.LogError(err)
		return nil, err
	}

	errChan, _ := app.Chain.ReserveBlockID(app.Self.NID)
	if err = <-errChan; err != nil {
		err = errors.Wrap(err, "SetValueHandler failed to ReserveBlockID")
		logger.LogError(err)
		return nil, err
	}

	errChan = app.Chain.AddNewBlock(block, app.Self.NID)
	if err = <-errChan; err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to AddNewBlock")
		logger.LogError(err)
		return nil, err
	}

	encGlobalKeyJSON, err := encGlobalKey.ToJSON()
	if err != nil {
		err = errors.Wrap(err, "AddVerifierNodeHandler failed to ToJSON")
		logger.LogError(err)
		return nil, err
	}

	protoMaster := ProtoNodeFromNode(app.Self)
	protoVerifier := ProtoNodeFromNode(verifier)

	resp := &NewNodeResponse{
		EncGlobalKey: encGlobalKeyJSON,
		Master:       protoMaster,
		Verifier:     protoVerifier,
	}

	return resp, nil
}
