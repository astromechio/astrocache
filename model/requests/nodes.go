package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	acrypto "github.com/astromechio/astrocache/crypto"
	"github.com/astromechio/astrocache/model"
)

// NewNodeRequest contains information for adding a new node
type NewNodeRequest struct {
	Address  string `json:"address"`
	NodeType string `json:"nodeType"`
	PubKey   []byte `json:"pubKey"`
	JoinCode string `json:"joinCode"`
}

// Path returns the path for a new node request
func (nr *NewNodeRequest) Path() string {
	typeSlug := ""

	switch nr.NodeType {
	case model.NodeTypeVerifier:
		typeSlug = "verifier"
	case model.NodeTypeWorker:
		typeSlug = "worker"
	}

	return fmt.Sprintf("v1/master/nodes/%s", typeSlug)
}

// FromRequest loads a new node request from an http request
func (nr *NewNodeRequest) FromRequest(r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(reqBody, nr)
}

// Verify verifies that the request is valid
func (nr *NewNodeRequest) Verify() error {
	if nr == nil {
		return errors.New("nr is nil")
	}

	if nr.Address == "" {
		return errors.New("nr.Address is empty")
	}

	if nr.NodeType != model.NodeTypeVerifier && nr.NodeType != model.NodeTypeWorker {
		return fmt.Errorf("nr.NodeType is %s, must be verifier or worker", nr.NodeType)
	}

	if len(nr.PubKey) == 0 {
		return errors.New("nr.PubKey length is 0")
	}

	if len(nr.JoinCode) == 0 {
		return errors.New("nr.JoinCode length is 0")
	}

	return nil
}

// NewNodeResponse contains everything a node needs to bootstrap istelf
type NewNodeResponse struct {
	*model.Node
	EncGlobalKey     *acrypto.Message `json:"encGlobalKey"`
	MasterPubKeyJSON []byte           `json:"masterPubKeyJSON"`
}
