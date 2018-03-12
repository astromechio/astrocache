package verifier

import (
	"net/http"

	"github.com/astromechio/astrocache/server"
	"github.com/astromechio/astrocache/server/verifier/handler"
	"github.com/gorilla/mux"
)

func router(config *server.Config) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods(http.MethodGet).Path("/v1/verifier/block/propose").HandlerFunc(handler.ProposeAddBlockHandler(config))
	mux.Methods(http.MethodGet).Path("/v1/verifier/block/check").HandlerFunc(handler.CheckBlockHandler(config))

	return mux
}
