package service

import (
	"context"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/logger"
	vservice "github.com/astromechio/astrocache/server/verifier/service"
	"github.com/pkg/errors"
)

// SetValueHandler handles setting values in the cache
func SetValueHandler(app *config.App, req *SetValueRequest) (*SetValueResponse, error) {
	verifier := app.NodeList.RandomVerifier()
	vReq := vReqFromWreq(req)

	conn, err := verifier.Dial()
	if err != nil {
		err = errors.Wrap(err, "SetValueHandler failed to Dial")
		logger.LogError(err)
		return nil, err
	}

	client := vservice.NewVerifierClient(conn)

	ctx := context.Background()
	_, err = client.SetValue(ctx, vReq)
	if err != nil {
		err = errors.Wrap(err, "SetValueHandler failed to SetValue")
		logger.LogError(err)
		return nil, err
	}

	return &SetValueResponse{}, nil
}

func vReqFromWreq(wReq *SetValueRequest) *vservice.SetValueRequest {
	return &vservice.SetValueRequest{
		Key:   wReq.Key,
		Value: wReq.Value,
	}
}
