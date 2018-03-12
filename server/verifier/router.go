package verifier

import (
	"net/http"

	"github.com/astromechio/astrocache/server"
	"github.com/astromechio/astrocache/server/verifier/handler"
	"github.com/gorilla/mux"
)

func router(config *server.Config) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods(http.MethodPost).Path("/v1/verifier/block/propose").HandlerFunc(handler.ProposeAddBlockHandler(config))

	// TODO: different method for check?
	mux.Methods(http.MethodPost).Path("/v1/verifier/block/check").HandlerFunc(handler.CheckBlockHandler(config))

	return mux
}
