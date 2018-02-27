package verifier

import (
	"github.com/astromechio/astrocache/modes/verifier/config"
	"github.com/gorilla/mux"
)

func router(config *config.Config) *mux.Router {
	mux := mux.NewRouter()

	// mux.Methods(http.MethodGet).Path("/v1/verifier/block/propose").HandlerFunc(common.CtxWrap(config.BaseConfig, handler.ProposeAddBlockHandler(chain)))

	return mux
}
