package master

import (
	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/server/master/handler"
	"github.com/gorilla/mux"
)

func router(app *config.App) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods("POST").Path("/v1/master/nodes/verifier").HandlerFunc(handler.AddVerifierNodeHandler(app))
	mux.Methods("GET").Path("/v1/master/chain").HandlerFunc(handler.GetEntireChainHandler(app))

	return mux
}
