package master

import (
	"net/http"

	"github.com/astromechio/astrocache/config"
	whandler "github.com/astromechio/astrocache/server/worker/handler"
	"github.com/gorilla/mux"
)

func router(app *config.App) *mux.Router {
	mux := mux.NewRouter()

	mux.Methods(http.MethodPost).Path("/v1/worker/block").HandlerFunc(whandler.AddBlockHandler(app))

	return mux
}
