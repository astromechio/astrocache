package service

import (
	"github.com/astromechio/astrocache/config"
	wservice "github.com/astromechio/astrocache/server/worker/service"
)

// AddBlockHandler handles adding blocks
func AddBlockHandler(app *config.App, req *AddBlockRequest) (*AddBlockResponse, error) {
	mReq := wReqFromMreq(req)

	_, err := wservice.AddBlockHandler(app, mReq)
	if err != nil {
		return nil, err
	}

	return &AddBlockResponse{}, nil
}

func wReqFromMreq(mReq *AddBlockRequest) *wservice.AddBlockRequest {
	return &wservice.AddBlockRequest{
		Block:        mReq.Block,
		ProposingNID: mReq.ProposingNID,
	}
}
