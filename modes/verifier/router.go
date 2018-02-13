package verifier

import (
	"net/http"

	"github.com/astromechio/astrocache/model/blockchain"
	"github.com/astromechio/astrocache/modes/common"
	"github.com/astromechio/astrocache/modes/verifier/handler"
	"github.com/gorilla/mux"
)

func router(config *Config, chain *blockchain.Chain) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods(http.MethodGet).Path("/v1/verifier/block/propose").HandlerFunc(common.CtxWrap(config.BaseConfig, handler.ProposeAddBlockHandler(chain)))

	return mux
}
