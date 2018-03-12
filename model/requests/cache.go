package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// KeyRequestKey and others are keys used for cache requests
const (
	KeyRequestKey = "key"
)

// SetValueRequest contains information for adding a new node
type SetValueRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Path returns the path for a new node request
func (sv *SetValueRequest) Path() string {
	return fmt.Sprintf("v1/value/%s", sv.Key)
}

// FromRequest loads a new node request from an http request
func (sv *SetValueRequest) FromRequest(r *http.Request) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(reqBody, sv); err != nil {
		return err
	}

	key := mux.Vars(r)[KeyRequestKey]
	if key == "" {
		return errors.New("No key found in request URL")
	}

	sv.Key = key

	return nil
}

// Verify verifies that the request is valid
func (sv *SetValueRequest) Verify() error {
	if sv == nil {
		return errors.New("sv is nil")
	}

	if sv.Key == "" {
		return errors.New("sv.Key is nil")
	}

	if sv.Value == "" {
		return errors.New("sv.Value is nil")
	}

	return nil
}
