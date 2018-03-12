package send

import (
	"fmt"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/transport"
)

// JoinNetwork requests that the current node be added to a network
func JoinNetwork(app *config.App, masterAddr, joinCode string) (*requests.NewNodeResponse, error) {
	request := &requests.NewNodeRequest{
		Node:     app.Self,
		JoinCode: joinCode,
	}

	masterAddrWithPort := fmt.Sprintf("%s:3000", masterAddr)

	url := transport.URLFromAddressAndPath(masterAddrWithPort, request.Path())

	resp := &requests.NewNodeResponse{}
	if err := transport.Post(url, request, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
