package requests

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/astromechio/astrocache/model/blockchain"
)

// ProposeBlockRequest contains information for adding a new node
type ProposeBlockRequest struct {
	Block    *blockchain.Block `json:"block"`
	MinerNID string            `json:"minerNid"`
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

	if pb.Block == nil {
		return errors.New("pb.Block is nil")
	}

	if pb.MinerNID == "" {
		return errors.New("pb.MinerNID is empty")
	}

	return nil
}

// CheckBlockRequest contains information for adding a new node
type CheckBlockRequest struct {
	Block *blockchain.Block `json:"block"`
}

// Path returns the path for a new node request
func (cb *CheckBlockRequest) Path() string {
	return "v1/verifier/block/check"
}

// FromRequest loads a new node request from an http request
func (cb *CheckBlockRequest) FromRequest(r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(reqBody, cb)
}

// Verify verifies that the request is valid
func (cb *CheckBlockRequest) Verify() error {
	if cb == nil {
		return errors.New("cb is nil")
	}

	if cb.Block == nil {
		return errors.New("cb.Block is nil")
	}

	return nil
}
