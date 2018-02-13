package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/pkg/errors"
)

// Post sends a POST request to a node with a request
func Post(node *model.Node, req requests.Request, res interface{}) error {
	reqURL := fmt.Sprintf("http://%s/%s", node.Address, req.Path())
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "Post failed to Marshal")
	}

	postRequest, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(reqJSON))
	if err != nil {
		return errors.Wrap(err, "Post failed to NewRequest")
	}

	response, err := http.DefaultClient.Do(postRequest)
	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "Post failed to ReadAll")
	}

	if err := json.Unmarshal(resBody, res); err != nil {
		return errors.Wrap(err, "Post failed to Unmarshal")
	}

	return nil
}

// Get sends a POST request to a node with a request
func Get(node *model.Node, req requests.Request, res interface{}) error {
	reqURL := fmt.Sprintf("http://%s/%s", node.Address, req.Path())

	getRequest, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return errors.Wrap(err, "Get failed to NewRequest")
	}

	response, err := http.DefaultClient.Do(getRequest)
	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "Get failed to ReadAll")
	}

	if err := json.Unmarshal(resBody, res); err != nil {
		return errors.Wrap(err, "Get failed to Unmarshal")
	}

	return nil
}
