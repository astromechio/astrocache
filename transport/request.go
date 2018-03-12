package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astromechio/astrocache/model/requests"
	"github.com/pkg/errors"
)

// Post sends a POST request to a node with a request
func Post(url string, req requests.Request, res interface{}) error {
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "Post failed to Marshal")
	}

	postRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqJSON))
	if err != nil {
		return errors.Wrap(err, "Post failed to NewRequest")
	}

	response, err := http.DefaultClient.Do(postRequest)
	if err != nil {
		return errors.Wrap(err, "Post failed to Do")
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Post (%s) returned non-200 status code %d", url, response.StatusCode)
	}

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "Post failed to ReadAll")
	}

	if res != nil {
		if err := json.Unmarshal(resBody, res); err != nil {
			return errors.Wrap(err, "Post failed to Unmarshal")
		}
	}

	return nil
}

// Get sends a POST request to a node with a request
func Get(url string, res interface{}) error {
	getRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "Get failed to NewRequest")
	}

	response, err := http.DefaultClient.Do(getRequest)
	if err != nil {
		return errors.Wrap(err, "Get failed to Do")
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("Post (%s) returned non-200 status code %d", url, response.StatusCode)
	}

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "Get failed to ReadAll")
	}

	if res != nil {
		if err := json.Unmarshal(resBody, res); err != nil {
			return errors.Wrap(err, "Get failed to Unmarshal")
		}
	}

	return nil
}

// URLFromAddressAndPath creates a URL from a root address and a path
func URLFromAddressAndPath(addr, path string) string {
	if !strings.HasPrefix(addr, "http://") {
		addr = fmt.Sprintf("http://%s", addr)
	}

	return fmt.Sprintf("%s/%s", addr, path)
}
