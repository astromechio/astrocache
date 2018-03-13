package verifier

import (
	"net/http"

	"github.com/astromechio/astrocache/config"
	"github.com/astromechio/astrocache/server/verifier/handler"
	"github.com/gorilla/mux"
)

func router(app *config.App) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods(http.MethodPost).Path("/v1/verifier/block/propose").HandlerFunc(handler.ProposeAddBlockHandler(app))

	// TODO: different method for check?
	mux.Methods(http.MethodPost).Path("/v1/verifier/block/check").HandlerFunc(handler.CheckBlockHandler(app))

	mux.Methods(http.MethodPost).Path("/v1/value/{key}").HandlerFunc(handler.SetValueHandler(app))

	return mux
}
