package service

import (
	"github.com/astromechio/astrocache/config"
	context "golang.org/x/net/context"
)

// VerifierService exposes a master service
type VerifierService struct {
	App *config.App
}

// ProposeBlock handles proposing blocks
func (vs *VerifierService) ProposeBlock(ctx context.Context, in *ProposeBlockRequest) (*ProposeBlockResponse, error) {
	return ProposeBlockHandler(vs.App, in)
}

// SetValue handles setting cache values
func (vs *VerifierService) SetValue(ctx context.Context, in *SetValueRequest) (*SetValueResponse, error) {
	return SetValueHandler(vs.App, in)
}
