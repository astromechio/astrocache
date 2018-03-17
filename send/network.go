package send

import (
	"context"

	"github.com/astromechio/astrocache/config"
	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/requests"
	mservice "github.com/astromechio/astrocache/server/master/service"
	"github.com/pkg/errors"
)

// JoinNetwork requests that the current node be added to a network
func JoinNetwork(app *config.App, master *model.Node, joinCode string) (*requests.NewNodeResponse, error) {
	protoSelf := mservice.ProtoNodeFromNode(app.Self)
	request := &mservice.NewNodeRequest{
		Node:     protoSelf,
		JoinCode: joinCode,
	}

	conn, err := master.Dial()
	if err != nil {
		return nil, errors.Wrap(err, "JoinNetwork failed to Dial")
	}

	client := mservice.NewMasterClient(conn)

	ctx := context.Background()
	resp, err := client.AddNode(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "JoinNetwork failed to AddNode")
	}

	encGlobalKey, err := acrypto.MessageFromJSON(resp.EncGlobalKey)
	if err != nil {
		return nil, errors.Wrap(err, "JoinNetwork failed to MessageFromJSON")
	}

	newMaster := mservice.NodeFromProtoNode(resp.Master)

	var verifier *model.Node
	if resp.Verifier != nil {
		verifier = mservice.NodeFromProtoNode(resp.Verifier)
	}

	response := &requests.NewNodeResponse{
		EncGlobalKey: encGlobalKey,
		Master:       newMaster,
		Verifier:     verifier,
		IsPrimary:    resp.IsPrimary,
	}

	return response, nil
}
