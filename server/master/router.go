package master

import (
	"github.com/astromechio/astrocache/server"
	"github.com/astromechio/astrocache/server/master/handler"
	"github.com/gorilla/mux"
)

func router(config *server.Config) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods("POST").Path("/v1/master/nodes/verifier").HandlerFunc(handler.AddVerifierNodeHandler(config))

	return mux
}
