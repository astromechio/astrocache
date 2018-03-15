package requests

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// ReserveIDRequest contains information for adding a new node
type ReserveIDRequest struct {
	ProposingNID string `json:"propNID"`
}

// Path returns the path for a new node request
func (rid *ReserveIDRequest) Path() string {
	return "v1/master/block/reserve"
}

// FromRequest loads a new node request from an http request
func (rid *ReserveIDRequest) FromRequest(r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.Unmarshal(reqBody, rid)
}

// Verify verifies that the request is valid
func (rid *ReserveIDRequest) Verify() error {
	if rid == nil {
		return errors.New("rid is nil")
	}

	if rid.ProposingNID == "" {
		return errors.New("rid.ProposingNID is empty")
	}

	return nil
}

// ReserveIDResponse is an ID reservation response
type ReserveIDResponse struct {
	BlockID string `json:"blockId"`
}
