package send

import (
	"context"

	"github.com/astromechio/astrocache/logger"
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/requests"
	vservice "github.com/astromechio/astrocache/server/verifier/service"
	wservice "github.com/astromechio/astrocache/server/worker/service"
	"github.com/pkg/errors"
)

// SetValueOnVerifier sends a value change request to a node
func SetValueOnVerifier(req *requests.SetValueRequest, node *model.Node) error {
	request := &vservice.SetValueRequest{
		Key:   req.Key,
		Value: req.Value,
	}

	conn, err := node.Dial()
	if err != nil {
		err = errors.Wrap(err, "SetValueOnVerifier failed to Dial")
		logger.LogError(err)
		return err
	}

	client := vservice.NewVerifierClient(conn)

	ctx := context.Background()
	_, err = client.SetValue(ctx, request)
	if err != nil {
		err = errors.Wrap(err, "SetValueOnVerifier failed to SetValue")
		logger.LogError(err)
		return err
	}

	return nil
}

// SetValueOnWorker sends a value change request to a node
func SetValueOnWorker(req *requests.SetValueRequest, node *model.Node) error {
	request := &wservice.SetValueRequest{
		Key:   req.Key,
		Value: req.Value,
	}

	conn, err := node.Dial()
	if err != nil {
		err = errors.Wrap(err, "SetValueOnWorker failed to Dial")
		logger.LogError(err)
		return err
	}

	client := wservice.NewWorkerClient(conn)

	ctx := context.Background()
	_, err = client.SetValue(ctx, request)
	if err != nil {
		err = errors.Wrap(err, "SetValueOnWorker failed to SetValue")
		logger.LogError(err)
		return err
	}

	return nil
}
