package requests

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	acrypto "github.com/astromechio/astrocache/crypto"
)

// ProposeBlockRequest contains information for adding a new node
type ProposeBlockRequest struct {
	TempID     string           `json:"tempID"`
	Data       *acrypto.Message `json:"data"`
	ActionType string           `json:"actionType"`
	PrevID     string           `json:"prevID"`
}

// Path returns the path for a new node request
func (pb *ProposeBlockRequest) Path() string {
	return "v1/verifier/block/propose"
}

// FromRequest loads a new node request from an http request
func (pb *ProposeBlockRequest) FromRequest(r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(reqBody, pb)
}

// Verify verifies that the request is valid
func (pb *ProposeBlockRequest) Verify() error {
	if pb == nil {
		return errors.New("pb is nil")
	}

	if pb.TempID == "" {
		return errors.New("pb.TempID is empty")
	}

	if pb.Data == nil {
		return errors.New("pb.Data is nil")
	}

	if pb.Data.KID == "" {
		return errors.New("pb.Data.KID is empty")
	}

	if pb.PrevID == "" {
		return errors.New("pb.PrevID is empty")
	}

	return nil
}

// ProposeBlockResponse contains everything a block needs to be committed
type ProposeBlockResponse struct {
	PrevHash   []byte `json:"prevHash"`
	IDMismatch bool   `json:"mismatch"`
}
