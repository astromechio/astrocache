package send

import (
	"github.com/astromechio/astrocache/model"
	"github.com/astromechio/astrocache/model/requests"
	"github.com/astromechio/astrocache/transport"
)

// SetValue sends a value change request to a node
func SetValue(req *requests.SetValueRequest, node *model.Node) error {
	url := transport.URLFromAddressAndPath(node.Address, req.Path())

	return transport.Post(url, req, nil)
}
