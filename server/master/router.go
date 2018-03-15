package master

import (
	"net/http"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/server/master/handler"
	whandler "github.com/astromechio/astrocache/server/worker/handler"
	"github.com/gorilla/mux"
)

func router(app *config.App) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods(http.MethodPost).Path("/v1/master/nodes/verifier").HandlerFunc(handler.AddVerifierNodeHandler(app))
	mux.Methods(http.MethodPost).Path("/v1/master/nodes/worker").HandlerFunc(handler.AddWorkerNodeHandler(app))

	mux.Methods(http.MethodGet).Path("/v1/master/chain").HandlerFunc(handler.GetEntireChainHandler(app))
	mux.Methods(http.MethodGet).Path("/v1/master/chain/after/{after}").HandlerFunc(handler.GetBlocksAfterHandler(app))

	mux.Methods(http.MethodPost).Path("/v1/worker/block").HandlerFunc(whandler.AddBlockHandler(app))

	return mux
}
