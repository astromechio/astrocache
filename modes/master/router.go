package master

import (
	"net/http"

	"github.com/astromechio/astrocache/modes/master/config"
	"github.com/astromechio/astrocache/modes/master/handler"
	"github.com/gorilla/mux"
)

type masterHandlerFunc func(*config.Config) http.HandlerFunc

func router(config *config.Config) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods("POST").Path("/v1/master/nodes/verifier").HandlerFunc(handler.AddVerifierNodeHandler(config))

	return mux
}
